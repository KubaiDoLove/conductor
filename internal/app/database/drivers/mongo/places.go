package mongo

import (
	"context"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlacesRepository struct {
	collection *mongo.Collection
}

func (p PlacesRepository) Create(ctx context.Context, place *models.Place) error {
	panic("implement me")
}

func (p PlacesRepository) Delete(ctx context.Context, placeID primitive.ObjectID) error {
	panic("implement me")
}

func (p PlacesRepository) ToBook(ctx context.Context, placeID primitive.ObjectID, booking *models.Booking) error {
	panic("implement me")
}

func (p PlacesRepository) CancelBooking(ctx context.Context, placeID primitive.ObjectID, bookingID primitive.ObjectID) error {
	panic("implement me")
}


