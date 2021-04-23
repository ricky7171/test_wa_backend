package chat

import (
	"github.com/ricky7171/test_wa_backend/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Reader interface
type Reader interface {
	FindByContact(contactId primitive.ObjectID, lastId string) ([]entity.Chat, error)
}

//Writer interface
type Writer interface {
	Create(message *entity.Message) (*entity.Message, error)
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	GetChatByContact(contactId, lastId string) ([]entity.Chat, error)
	MakePreparedMessage(currentUserId, currentUserName, receiverPhone, message string) (*entity.Message, error)
	CreateMessage(message *entity.Message) (*entity.Message, error)
}
