package entity

import (
	"time"

	"github.com/ricky7171/test_wa_backend/internal/failure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserWithName struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type Contact struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Users     []UserWithName     `json:"users_info" bson:"users"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
}

func NewContact(firstUserId, secondUserId primitive.ObjectID, firstUserName, secondUserName string) (*Contact, error) {
	//1. make new entity
	firstUser := UserWithName{
		ID:   firstUserId,
		Name: firstUserName,
	}
	secondUser := UserWithName{
		ID:   secondUserId,
		Name: secondUserName,
	}
	newContact := &Contact{
		ID:        primitive.NewObjectID(),
		Users:     []UserWithName{firstUser, secondUser},
		CreatedAt: time.Now(),
	}

	//2. validate that entity
	if err := newContact.Validate(); err != nil {
		return nil, err
	}

	//3. return new entity
	return newContact, nil
}

func (c *Contact) Validate(fields ...string) error {
	if len(fields) == 0 {
		fields = []string{"users"}
	}
	for _, value := range fields {
		if value == "users" {
			if c.Users[0].ID.Hex() == "" {
				return failure.ErrFieldRequired("first user id")
			}
			if c.Users[1].ID.Hex() == "" {
				return failure.ErrFieldRequired("second user id")
			}
			if c.Users[0].Name == "" {
				return failure.ErrFieldRequired("first user name")
			}
			if c.Users[1].Name == "" {
				return failure.ErrFieldRequired("second user name")
			}
			if len(c.Users[0].Name) < 2 || len(c.Users[0].Name) > 100 {
				return failure.ErrFieldLenConstraint("first user name", "2", "100")
			}
			if len(c.Users[1].Name) < 2 || len(c.Users[1].Name) > 100 {
				return failure.ErrFieldLenConstraint("second user name", "2", "100")
			}
		}
	}
	return nil
}
