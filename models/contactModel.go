package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	Users     []primitive.ObjectID `json:"users"`
	CreatedAt time.Time            `json:"created_at"`
}

type UserWithName struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type ContactWithName struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Users     []UserWithName     `json:"users_info" bson:"users_info"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
}
