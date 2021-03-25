package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Phone        string             `json:"phone" bson:"phone" validate:"required"`
	Password     string             `json:",omitempty" bson:"password" validate:"required,min=6"`
	Token        string             `json:"token" bson:",omitempty"`
	RefreshToken string             `json:"refreshToken" bson:",omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}
