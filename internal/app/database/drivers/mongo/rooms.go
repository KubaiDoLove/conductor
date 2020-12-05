package mongo

import (
	"context"
	"github.com/KubaiDoLove/conductor/internal/app/database/drivers"
	"github.com/KubaiDoLove/conductor/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomsRepository struct {
	collection *mongo.Collection
}

func (r RoomsRepository) Create(ctx context.Context, room *models.Room) error {
	if room == nil {
		return drivers.ErrEmptyRoom
	}

	room.ID = primitive.NewObjectID()
	if _, err := r.collection.InsertOne(ctx, room); err != nil {
		return err
	}

	return nil
}

func (r RoomsRepository) RoomByID(ctx context.Context, roomID primitive.ObjectID) (*models.Room, error) {
	filter := bson.D{{"_id", roomID}}

	room := new(models.Room)
	err := r.collection.FindOne(ctx, filter, nil).Decode(room)

	switch err {
	case nil:
		return room, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrRoomDoesNotExist
	default:
		return nil, err
	}
}

func (r RoomsRepository) Update(ctx context.Context, room *models.Room) error {
	if room == nil {
		return drivers.ErrEmptyRoom
	}

	filter := bson.D{{Key: "_id", Value: room.ID}}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: room.Name},
				{Key: "maxLandingPercent", Value: room.MaxLandingPercent},
				{Key: "socialDistance", Value: room.SocialDistance},
			},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrRoomDoesNotExist
	}

	return nil
}

func (r RoomsRepository) Delete(ctx context.Context, roomID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: roomID}})
	return err
}

