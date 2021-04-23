package entity

import (
	"time"

	"github.com/ricky7171/test_wa_backend/internal/failure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	ContactId primitive.ObjectID `json:"contact_id" bson:"contactId"`
	SenderId  primitive.ObjectID `json:"sender_id" bson:"senderId"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
}

func NewChat(contactId primitive.ObjectID, senderId primitive.ObjectID, message string) (*Chat, error) {
	//1. make new entity
	newChat := &Chat{
		ID:        primitive.NewObjectID(),
		ContactId: contactId,
		SenderId:  senderId,
		Message:   message,
		CreatedAt: time.Now(),
	}

	//2. validate that entity
	if err := newChat.Validate(); err != nil {
		return nil, err
	}

	//3. return new entity
	return newChat, nil
}

func (c *Chat) Validate(fields ...string) error {
	if len(fields) == 0 {
		fields = []string{"contactId", "senderId", "message"}
	}
	for _, field := range fields {
		if field == "contactId" {
			if c.ContactId.Hex() == "" {
				return failure.ErrFieldRequired(field)
			}
		} else if field == "senderId" {
			if c.SenderId.Hex() == "" {
				return failure.ErrFieldRequired(field)
			}
		} else if field == "message" {
			if c.Message == "" {
				return failure.ErrFieldRequired(field)
			}
			if len(c.Message) < 2 || len(c.Message) > 100 {
				return failure.ErrFieldLenConstraint(field, "2", "100")
			}
		}
	}
	return nil
}
