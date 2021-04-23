package entity

import (
	"errors"
	"time"

	"github.com/ricky7171/test_wa_backend/internal/failure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Phone        string             `json:"phone" bson:"phone"`
	Password     string             `json:",omitempty" bson:"password"`
	Token        string             `json:"token,omitempty" bson:",omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:",omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updatedAt"`
}

func NewUser(name, phone, password string) (*User, error) {
	//1. make new entity
	newUser := &User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Phone:     phone,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//2. validate that entity
	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	//3. hash password
	err := newUser.HashPassword()
	if err != nil {
		return nil, failure.ErrHashPassword()
	}

	//4. return new entity
	return newUser, nil
}

func (u *User) Validate(fields ...string) error {
	if len(fields) == 0 {
		fields = []string{"name", "phone", "password"}
	}
	for _, field := range fields {
		if field == "name" {
			if u.Name == "" {
				return failure.ErrFieldRequired(field)
			}
			if len(u.Name) < 2 || len(u.Name) > 100 {
				return failure.ErrFieldLenConstraint(field, "2", "100")
			}
		} else if field == "phone" {
			if u.Phone == "" {
				return failure.ErrFieldRequired(field)
			}
			if len(u.Phone) < 2 || len(u.Phone) > 30 {
				return failure.ErrFieldLenConstraint(field, "2", "30")
			}
		} else if field == "password" {
			if u.Password == "" {
				return failure.ErrFieldRequired(field)
			}
			if len(u.Password) < 6 {
				return failure.ErrFieldLenConstraint(field, "6", "")
			}
		}
	}
	return nil
}

//HashPassword is used to encrypt the password before it is stored in the DB
func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

//used to verify whether this user password is match with provided password (plain password)
func (u *User) VerifyPassword(plainPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))

	if err != nil {
		return errors.New("password doesn't match")
	}

	return nil
}

//used to set token field on this object
func (u *User) SetToken(token, refreshToken string) {
	u.Token = token
	u.RefreshToken = refreshToken
}
