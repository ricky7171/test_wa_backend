package entity

import (
	"time"

	"github.com/ricky7171/test_wa_backend/internal/failure"
)

type Message struct {
	Data       string    `json:"data"`
	FromUserId string    `json:"from_user_id"`
	ToUserId   string    `json:"to_user_id"`
	CreatedAt  time.Time `json:",omitempty"`
	ContactId  string    `json:"contact_id"`
}

func NewMessage(data, fromUserId, ToUserId, contactId string) (*Message, error) {
	//1. make new entity
	newMessage := &Message{
		Data:       data,
		FromUserId: fromUserId,
		ToUserId:   ToUserId,
		ContactId:  contactId,
		CreatedAt:  time.Now(),
	}

	//2. validate that entity
	if err := newMessage.Validate(); err != nil {
		return nil, err
	}

	//4. return new entity
	return newMessage, nil
}

func (m *Message) Validate(fields ...string) error {
	if len(fields) == 0 {
		fields = []string{"data", "fromUserId", "toUserId", "contactId"}
	}
	for _, field := range fields {
		if field == "data" {
			if m.Data == "" {
				return failure.ErrFieldRequired(field)
			}
			if len(m.Data) < 2 || len(m.Data) > 100 {
				return failure.ErrFieldLenConstraint(field, "2", "100")
			}
		} else if field == "fromUserId" {
			if m.FromUserId == "" {
				return failure.ErrFieldRequired(field)
			}
		} else if field == "toUserId" {
			if m.ToUserId == "" {
				return failure.ErrFieldRequired(field)
			}
		} else if field == "contactId" {
			if m.ContactId == "" {
				return failure.ErrFieldRequired(field)
			}
		}
	}
	return nil
}
