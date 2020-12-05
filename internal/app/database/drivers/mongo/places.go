package mongo

import (
	"context"
	"github.com/KubaiDoLove/conductor/internal/app/database/drivers"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (p PlacesRepository) Delete(ctx context.Context, roomID primitive.ObjectID, placeID primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: roomID}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{{Key: "places._id", Value: placeID}}},
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

func (p PlacesRepository) ToBook(ctx context.Context, placeID primitive.ObjectID, booking *models.Booking) error {
	return nil
}

func (p PlacesRepository) CancelBooking(ctx context.Context, placeID primitive.ObjectID, bookingID primitive.ObjectID) error {
	return nil
}
