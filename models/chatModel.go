package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID         primitive.ObjectID `bson:"_id"`
	Contact_id primitive.ObjectID `json:"contact_id"`
	Sender_id  primitive.ObjectID `json:"sender_id"`
	Message    string             `json:"message"`
	Created_at time.Time          `json:"created_at"`
}
