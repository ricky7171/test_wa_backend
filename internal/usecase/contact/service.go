package contact

import (
	"github.com/ricky7171/test_wa_backend/internal/entity"
	"github.com/ricky7171/test_wa_backend/internal/failure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ContactRepo Repository
}

//create new service
func NewService(r Repository) *Service {
	return &Service{
		ContactRepo: r,
	}
}

//used to get Contact
func (s *Service) GetContactByUser(userId string, userName string) ([]entity.Contact, error) {
	//note : this function need userId and userName because in contact collection it save not only user id but also user name (to optimize performance when get contact)

	//1. convert contactId to contactObjectId
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, failure.ErrConvertObjectIdToHex()
	}

	//2. find contact by user
	contacts, err := s.ContactRepo.FindByUser(&entity.UserWithName{ID: userObjectId, Name: userName}, nil)
	if err != nil {
		return nil, err
	}

	//3. validate entity
	for _, contact := range contacts {
		if err := contact.Validate(); err != nil {
			return nil, err
		}
	}

	return contacts, nil

}
