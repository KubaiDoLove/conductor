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

	filter := bson.D{{Key: "places._id", Value: placeID}}
	room := new(models.Room)
	if err := p.collection.FindOne(ctx, filter).Decode(room); err != nil {
		if err == mongo.ErrNoDocuments {
			return drivers.ErrPlaceDoesNotExist
		}
		return err
	}

	now := time.Now()
	busyPlaces := make([]models.Place, 0, len(room.Places))
	placeToUpdateIdx := -1
	for placeIdx, place := range room.Places {
		if place.ID.Hex() == placeID.Hex() {
			placeToUpdateIdx = placeIdx
		}

		for _, b := range place.Bookings {
			placeIsBusy := b.StartTime.Before(now) || b.StartTime.Equal(now)
			placeWillBeBusy := b.EndTime.After(now)

			if placeIsBusy && placeWillBeBusy {
				busyPlaces = append(busyPlaces, place)
				break
			}
		}
	}

	placeOpenHour, placeCloseHour := room.Places[placeToUpdateIdx].WorkingHours[0], room.Places[placeToUpdateIdx].WorkingHours[1]
	bookingBeforeOpen := booking.StartTime.Hour() < placeOpenHour
	bookingAfterClose := booking.EndTime.Hour() > placeCloseHour
	if placeOpenHour == 0 {
		bookingBeforeOpen = false
	}
	if placeCloseHour == 0 {
		bookingAfterClose = false
	}
	if bookingBeforeOpen || bookingAfterClose {
		return fmt.Errorf("place works between %d:00-%d:00", placeOpenHour, placeCloseHour)
	}

	freePlaceWillBeBooked := true
	for _, busyPlace := range busyPlaces {
		if busyPlace.ID.Hex() == placeID.Hex() {
			freePlaceWillBeBooked = false
			break
		}
	}

	if freePlaceWillBeBooked {
		nextBusyPlacesCount := len(busyPlaces) + 1
		nextLandingPercent := nextBusyPlacesCount * 100 / len(room.Places)
		if nextLandingPercent > int(room.MaxLandingPercent) {
			return fmt.Errorf("it is impossible to book more rooms, max landing percent is %d", room.MaxLandingPercent)
		}

		isNotAllowedBySocialDistance := false

		for idx, p := range room.Places {
			if idx == placeToUpdateIdx {
				continue
			}

			conflictByX := math.AbsInt(p.X-room.Places[placeToUpdateIdx].X) < int(room.SocialDistance)
			conflictByY := math.AbsInt(p.Y-room.Places[placeToUpdateIdx].Y) < int(room.SocialDistance)
			if conflictByX && conflictByY {
				isNotAllowedBySocialDistance = true
				break
			}
		}

		if isNotAllowedBySocialDistance {
			return fmt.Errorf("cannot book this place: minimal social distance is %d", room.SocialDistance)
		}
	}

	dayBookings := make([]models.Booking, 0)
	for _, b := range room.Places[placeToUpdateIdx].Bookings {
		bYear, bMonth, bDay := b.StartTime.Date()
		toBYear, toBMonth, toBDay := booking.StartTime.Date()

		if bYear == toBYear && bMonth == toBMonth && bDay == toBDay {
			dayBookings = append(dayBookings, b)
		}
	}

	for _, b := range dayBookings {
		timeConflict := booking.EndTime.Before(b.EndTime) && booking.EndTime.After(b.StartTime)
		if timeConflict {
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

func (p PlacesRepository) CancelBooking(ctx context.Context, placeID primitive.ObjectID, bookingID primitive.ObjectID) error {
	return nil
}
