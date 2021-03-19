package models

import (
	"time"
)

//contoh message
// {data : "blablablabla", fromUserId : 1, toUserId : 2}
// ini masih belum tau buat apa, tapi intinya ini adalah data antara message dengan room
type Message struct {
	Data       string    `json:"data"`
	FromUserId string    `json:"fromUserId"`
	ToUserId   string    `json:"toUserId"`
	Created_at time.Time `json:",omitempty"`
	Contact_id string    `json:"contact_id"`
}

type NewChat struct {
	Phone      string `json:"phone" validate:"min=2,max=100"`
	Message    string `json:"message" validate:"required,min=2,max=100"`
	Contact_id string `json:"contact_id" validate:"min=2,max=100"`
}
