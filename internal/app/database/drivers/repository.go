package drivers

import (
	"context"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomsRepository interface {
	All(ctx context.Context) ([]models.Room, error)
	Create(ctx context.Context, room *models.Room) error
	RoomByID(ctx context.Context, roomID primitive.ObjectID) (*models.Room, error)
	Update(ctx context.Context, room *models.Room) error
	Delete(ctx context.Context, roomID primitive.ObjectID) error
}

type PlacesRepository interface {
	Create(ctx context.Context, roomID primitive.ObjectID, place *models.Place) error
	Delete(ctx context.Context, placeID primitive.ObjectID) error
	AddBooking(ctx context.Context, placeID primitive.ObjectID, booking *models.Booking) error
	CancelBooking(ctx context.Context, bookingID primitive.ObjectID) error
}
