package models

import (
	"time"
)

//contoh message
// {data : "blablablabla", fromUserId : 1, toUserId : 2}
// ini masih belum tau buat apa, tapi intinya ini adalah data antara message dengan room
type Message struct {
	Data       string    `json:"data"`
	FromUserId string    `json:"from_user_id"`
	ToUserId   string    `json:"to_user_id"`
	CreatedAt  time.Time `json:",omitempty"`
	ContactId  string    `json:"contact_id"`
}

type NewChat struct {
	Phone     string `json:"phone" validate:"min=2,max=100"`
	Message   string `json:"message" validate:"required,min=2,max=100"`
	ContactId string `json:"contact_id" validate:"min=2,max=100"`
}
