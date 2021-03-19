package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID         primitive.ObjectID   `bson:"_id"`
	Users      []primitive.ObjectID `json:"users"`
	Created_at time.Time            `json:"created_at"`
}
