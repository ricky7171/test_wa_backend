package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Phone        string             `json:"phone" bson:"phone" validate:"required"`
	Password     string             `json:",omitempty" bson:"password" validate:"required,min=6"`
	Token        string             `json:"token,omitempty" bson:",omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:",omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updatedAt"`
}
