package models

import (
	"time"
)

type Message struct {
	Data       string    `json:"data"`
	FromUserId string    `json:"from_user_id"`
	ToUserId   string    `json:"to_user_id"`
	CreatedAt  time.Time `json:",omitempty"`
	ContactId  string    `json:"contact_id"`
}

type NewChat struct {
	Phone     string `json:"phone"`
	Message   string `json:"message" validate:"required,min=2,max=100"`
	ContactId string `json:"contact_id"`
}
