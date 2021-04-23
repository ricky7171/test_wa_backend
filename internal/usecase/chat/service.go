package chat

import (
	"github.com/ricky7171/test_wa_backend/internal/entity"
	"github.com/ricky7171/test_wa_backend/internal/failure"
	"github.com/ricky7171/test_wa_backend/internal/usecase/contact"
	"github.com/ricky7171/test_wa_backend/internal/usecase/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	chatRepo    Repository
	userRepo    user.Repository
	contactRepo contact.Repository
}

//create new service
func NewService(chatRepo Repository, userRepo user.Repository, contactRepo contact.Repository) *Service {
	return &Service{
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		contactRepo: contactRepo,
	}
}

//used to get chat by contact
func (s *Service) GetChatByContact(contactId, lastId string) ([]entity.Chat, error) {
	//1. convert contactId to contactObjectId
	contactObjectId, err := primitive.ObjectIDFromHex(contactId)
	if err != nil {
		return nil, failure.ErrConvertObjectIdToHex()
	}

	//2. find chat by contact
	chats, err := s.chatRepo.FindByContact(contactObjectId, lastId)
	if err != nil {
		return nil, err
	}

	//3. validate entities
	for _, chat := range chats {
		if err := chat.Validate(); err != nil {
			return nil, err
		}
	}

	return chats, nil

}

//used to make message entity
func (s *Service) MakePreparedMessage(currentUserId, currentUserName, receiverPhone, message string) (*entity.Message, error) {

	//1. find user by receiver phone
	receiverUser, err := s.userRepo.FindByPhone(receiverPhone)
	if err != nil {
		return nil, err
	}

	//2. validate entity
	if err := receiverUser.Validate(); err != nil {
		return nil, err
	}

	//3. get receiver user id in string
	toUserId := receiverUser.ID.Hex()

	//4. if user found, check whether they have contact each other or not, But before that make current UserWithName first
	currentUserObjectId, err := primitive.ObjectIDFromHex(currentUserId)
	if err != nil {
		return nil, failure.ErrConvertObjectIdToHex()
	}
	currentUserWithName := &entity.UserWithName{ID: currentUserObjectId, Name: currentUserName}

	//5. check whether current user have contact each other or not
	var contactId string
	contactExists, err := s.contactRepo.FindByUser(
		currentUserWithName,
		&entity.UserWithName{ID: receiverUser.ID, Name: receiverUser.Name},
	)

	//6. validate entity
	for _, contact := range contactExists {
		if err := contact.Validate(); err != nil {
			return nil, err
		}
	}

	//7. if contact doesn't found, then create contact first
	if err != nil {
		newContact, err := entity.NewContact(currentUserWithName.ID, receiverUser.ID, currentUserWithName.Name, receiverUser.Name)
		if err != nil {
			return nil, err
		}
		_, err = s.contactRepo.Create(newContact)
		if err != nil {
			return nil, err
		}
		contactId = newContact.ID.Hex()
	} else {
		contactId = contactExists[0].ID.Hex()
	}

	//8. return message
	//8.a. make message object first
	messageEntity, err := entity.NewMessage(message, currentUserId, toUserId, contactId)
	if err != nil {
		return nil, err
	}

	//8.b. return it
	return messageEntity, nil
}

func (s *Service) CreateMessage(message *entity.Message) (*entity.Message, error) {
	return s.chatRepo.Create(message)
}
