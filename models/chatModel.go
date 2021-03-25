package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `bson:"_id"`
	ContactId primitive.ObjectID `json:"contactId" bson:"contactId"`
	SenderId  primitive.ObjectID `json:"senderId" bson:"senderId"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
