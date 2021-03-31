package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserWithName struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type ContactWithName struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Users     []UserWithName     `json:"users_info" bson:"users"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
}
