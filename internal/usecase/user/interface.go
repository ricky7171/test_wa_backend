package user

import (
	"github.com/ricky7171/test_wa_backend/internal/entity"
)

//Reader interface
type Reader interface {
	CheckUserExists(phone string) (bool, error)
	FindByPhone(phone string) (*entity.User, error)
	FindById(userId string) (*entity.User, error)
}

//Writer interface
type Writer interface {
	Create(e *entity.User) (string, error)
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	CreateUser(name, phone, password string) (string, error)
	Authenticate(phone, password string) (*entity.User, error)
	RefreshToken(userId string) (*entity.User, error)
}
