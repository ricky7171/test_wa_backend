package contact

import (
	"github.com/ricky7171/test_wa_backend/internal/entity"
)

//Reader interface
type Reader interface {
	FindByUser(firstUser *entity.UserWithName, secondUser *entity.UserWithName) ([]entity.Contact, error)
}

//Writer interface
type Writer interface {
	Create(contact *entity.Contact) (string, error)
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	GetContactByUser(userId string, userName string) ([]entity.Contact, error)
}
