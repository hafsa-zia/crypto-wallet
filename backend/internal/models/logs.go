package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Event     string             `bson:"event" json:"event"`
	Details   string             `bson:"details" json:"details"`
	IP        string             `bson:"ip,omitempty" json:"ip,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
