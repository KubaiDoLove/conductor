package mongo

import (
	"context"
	"fmt"
	"github.com/KubaiDoLove/conductor/internal/app/database/drivers"
	"github.com/KubaiDoLove/conductor/internal/app/math"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type PlacesRepository struct {
	collection *mongo.Collection
}

func (p PlacesRepository) Create(ctx context.Context, roomID primitive.ObjectID, place *models.Place) error {
	if place == nil {
		return drivers.ErrEmptyPlace
	}

	filter := bson.D{{Key: "_id", Value: roomID}}
	room := new(models.Room)

	if err := p.collection.FindOne(ctx, filter).Decode(room); err != nil {
		if err == mongo.ErrNoDocuments {
			return drivers.ErrRoomDoesNotExist
		}
		return err
	}

	isInvalidX := place.X < 0 || place.X > room.Size[0]
	isInvalidY := place.Y < 0 || place.Y > room.Size[1]
	if isInvalidX || isInvalidY {
		return drivers.ErrInvalidPlace
	}

	for _, p := range room.Places {
		if p.X == place.X && p.Y == place.Y {
			return drivers.ErrPlaceTaken
		}
	}

	place.ID = primitive.NewObjectID()
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{{Key: "places", Value: place}}},
	}

	result, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrRoomDoesNotExist
	}

	return nil
}

func (p PlacesRepository) Delete(ctx context.Context, placeID primitive.ObjectID) error {
	filter := bson.D{{Key: "places._id", Value: placeID}}
	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{
					Key: "places",
					Value: bson.D{
						{
							Key:   "_id",
							Value: placeID,
						},
					},
				},
			},
		},
	}

	result, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrPlaceDoesNotExist
	}

	return nil
}

func (p PlacesRepository) AddBooking(ctx context.Context, placeID primitive.ObjectID, booking *models.Booking) error {
	if booking == nil {
		return drivers.ErrEmptyBooking
	}

	filter := bson.D{{Key: "places._id", Value: placeID}} // Ищем комнату по placeID
	room := new(models.Room)
	if err := p.collection.FindOne(ctx, filter).Decode(room); err != nil {
		if err == mongo.ErrNoDocuments {
			return drivers.ErrPlaceDoesNotExist
		}
		return err
	}

	now := time.Now().Local()
	busyPlaces := make([]models.Place, 0, len(room.Places)) // ищем уже занятые места
	placeToUpdateIdx := -1
	for placeIdx, place := range room.Places {
		if place.ID.Hex() == placeID.Hex() {
			placeToUpdateIdx = placeIdx // сохраняем индекс нашего места для дальнейших манипуляций
		}

		for _, b := range place.Bookings {
			placeIsBusy := b.StartTime.Local().Before(now) || b.StartTime.Local().Equal(now)
			placeWillBeBusy := b.EndTime.Local().After(now)

			if placeIsBusy && placeWillBeBusy { //  Если место занято и будет занято еще, то записываем в занятые места
				busyPlaces = append(busyPlaces, place)
				break
			}
		}
	}

	// Проверяем на валидность бронирования в рабочие часы
	openHour, closeHour := room.WorkingHours[0], room.WorkingHours[1]
	bookingBeforeOpen := booking.StartTime.Local().Hour() < openHour
	bookingAfterClose := booking.EndTime.Local().Hour() > closeHour
	if openHour == 0 {
		bookingBeforeOpen = false
	}
	if closeHour == 0 {
		bookingAfterClose = false
	}
	if bookingBeforeOpen || bookingAfterClose {
		return fmt.Errorf("room works between %d:00-%d:00", openHour, closeHour)
	}

	freePlaceWillBeBooked := true // Скорее всего не занятое место будет забронировано
	for _, busyPlace := range busyPlaces {
		if busyPlace.ID.Hex() == placeID.Hex() {
			freePlaceWillBeBooked = false // Но и уже занятое место на будущее могут забронировать
			break
		}
	}

	if freePlaceWillBeBooked {
		nextBusyPlacesCount := len(busyPlaces) + 1
		nextLandingPercent := nextBusyPlacesCount * 100 / len(room.Places)
		if nextLandingPercent > int(room.MaxLandingPercent) {
			// Если кол-во занятых мест + 1 превышает максимально разрешенный, то дропаем бронь
			return fmt.Errorf("it is impossible to book more rooms, max landing percent is %d", room.MaxLandingPercent)
		}
	}

	isNotAllowedBySocialDistance := false // Предполагаем, что нам разрешено забронировать новое место по правилу соц.дистанции
	if room.SocialDistance > 0 {          // Если нужно разделять по соц.дистанции
		for idx, p := range busyPlaces {
			if idx == placeToUpdateIdx { // Это наше же место, проверять на расстояние не надо
				continue
			}

			// Разницы координат по модулю
			conflictByX := math.AbsInt(p.X-room.Places[placeToUpdateIdx].X) <= int(room.SocialDistance)
			conflictByY := math.AbsInt(p.Y-room.Places[placeToUpdateIdx].Y) <= int(room.SocialDistance)
			if conflictByX && conflictByY {
				isNotAllowedBySocialDistance = true
				break
			}
		}
	}
	if isNotAllowedBySocialDistance {
		return fmt.Errorf("cannot book this place: minimal social distance is %d", room.SocialDistance)
	}

	dayBookings := make([]models.Booking, 0) // Бронирования по тому же дню
	for _, b := range room.Places[placeToUpdateIdx].Bookings {
		bYear, bMonth, bDay := b.StartTime.Local().Date()
		toBYear, toBMonth, toBDay := booking.StartTime.Local().Date()

		if bYear == toBYear && bMonth == toBMonth && bDay == toBDay {
			dayBookings = append(dayBookings, b)
		}
	}

	for _, b := range dayBookings {
		// Если потенциальное время бронирования входит в промежуток одного из бронирований в тот же день, дропем бронь
		startTimeConflict := booking.StartTime.Local().Before(b.EndTime.Local()) && booking.StartTime.Local().After(b.StartTime.Local())
		endTimeConflict := booking.EndTime.Local().Before(b.EndTime.Local()) && booking.EndTime.Local().After(b.StartTime.Local())
		equalTimeConflict := booking.EndTime.Local().Equal(b.EndTime.Local()) || booking.StartTime.Local().Equal(b.StartTime.Local())
		if startTimeConflict || endTimeConflict || equalTimeConflict {
			return drivers.ErrBookingTimeConflict
		}
	}

	booking.ID = primitive.NewObjectID()
	room.Places[placeToUpdateIdx].Bookings = append(room.Places[placeToUpdateIdx].Bookings, *booking)

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "places", Value: room.Places},
			},
		},
	}
	result, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrPlaceDoesNotExist
	}

	return nil
}

func (p PlacesRepository) CancelBooking(ctx context.Context, bookingID primitive.ObjectID) error {
	filter := bson.D{{Key: "places.bookings._id", Value: bookingID}}
	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{
					Key: "places.$[].bookings",
					Value: bson.D{
						{
							Key:   "_id",
							Value: bookingID,
						},
					},
				},
			},
		},
	}

	result, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrBookingDoesNotExist
	}

	return nil
}

func (p PlacesRepository) Suggestions(ctx context.Context, placeID primitive.ObjectID) ([]models.OpenCloseHours, error) {
	filter := bson.D{{Key: "places._id", Value: placeID}} // Ищем комнату по placeID
	room := new(models.Room)
	if err := p.collection.FindOne(ctx, filter).Decode(room); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, drivers.ErrPlaceDoesNotExist
		}
		return nil, err
	}

	now := time.Now().Local()
	earliestHour := now.Hour()
	latestHour := room.WorkingHours[1]
	if room.WorkingHours[0] > earliestHour || room.WorkingHours[0] == 0 {
		earliestHour = room.WorkingHours[0]
	}

	currentPlace := new(models.Place)
	for _, place := range room.Places {
		if place.ID.Hex() == placeID.Hex() {
			currentPlace = &place
			break
		}
	}

	skipHours := make([]models.OpenCloseHours, 0)
	for _, booking := range currentPlace.Bookings {
		for hour := booking.StartTime.Local().Hour(); hour < booking.EndTime.Local().Hour(); hour++ {
			skipHours = append(skipHours, models.OpenCloseHours{hour, hour + 1})
		}
	}

	freeTime := make([]models.OpenCloseHours, 0)
	for hour := earliestHour; hour < latestHour; hour++ {
		shouldAddTime := true

		for _, hoursPair := range skipHours {
			if hour == hoursPair[0] && hour+1 == hoursPair[1] {
				shouldAddTime = false
				break
			}
		}

		if shouldAddTime {
			freeTime = append(freeTime, models.OpenCloseHours{hour, hour + 1})
		}
	}

	return freeTime, nil
}
