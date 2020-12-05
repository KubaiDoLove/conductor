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

type RoomSize [2]int

type Room struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Size              RoomSize           `bson:"roomSize" json:"roomSize"`
	Places            []Place            `bson:"places" json:"places"`
	MaxLandingPercent uint8              `bson:"maxLandingPercent" json:"maxLandingPercent"`
	SocialDistance    uint               `bson:"socialDistance" json:"socialDistance"`
}
