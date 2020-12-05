package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	StartTime time.Time          `bson:"startTime" json:"startTime"`
	EndTime   time.Time          `bson:"endTime" json:"endTime"`
}

type Place struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	X        int                `bson:"x" json:"x"`
	Y        int                `bson:"y" json:"y"`
	Bookings []Booking          `bson:"bookings" json:"bookings"`
}

func NewPlace(x, y int) *Place {
	return &Place{
		X:        x,
		Y:        y,
		Bookings: make([]Booking, 0),
	}
}

// Размер - условные X и Y, только ширина и длина
type RoomSize [2]int

type Room struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Size              RoomSize           `bson:"roomSize" json:"roomSize"`
	Places            []Place            `bson:"places" json:"places"`
	MaxLandingPercent uint8              `bson:"maxLandingPercent" json:"maxLandingPercent"`
	SocialDistance    uint               `bson:"socialDistance" json:"socialDistance"`
}

func NewRoom(name string, size RoomSize, maxLandPercent uint8, socialDistance uint) *Room {
	if maxLandPercent > 100 {
		maxLandPercent = 100
	}

	return &Room{
		Name:              name,
		Size:              size,
		Places:            make([]Place, 0),
		MaxLandingPercent: maxLandPercent,
		SocialDistance:    socialDistance,
	}
}
