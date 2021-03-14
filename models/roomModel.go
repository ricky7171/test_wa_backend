package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Room is the model that governs all notes objects retrived or inserted into the DB
type Room struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Chat_history []ChatHistory        `json:"chat_history"`
	Users        []primitive.ObjectID `json:"users"`
	Created_at   time.Time            `json:"created_at"`
	Updated_at   time.Time            `json:"updated_at"`
	Room_id      string               `json:"room_id"`
}

type ChatHistory struct {
	User_id    string    `json:"user_id"`
	Message    string    `json:"message"`
	Created_at time.Time `json:"created_at"`
}
