package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	ContactId primitive.ObjectID `json:"contact_id" bson:"contactId"`
	SenderId  primitive.ObjectID `json:"sender_id" bson:"senderId"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
}
