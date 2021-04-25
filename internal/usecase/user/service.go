package user

import (
	"github.com/ricky7171/test_wa_backend/internal/entity"
	"github.com/ricky7171/test_wa_backend/internal/failure"
	"github.com/ricky7171/test_wa_backend/internal/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	userRepo    Repository
	tokenHelper helper.Token
}

//create new service
func NewService(r Repository, t helper.Token) *Service {
	return &Service{
		userRepo:    r,
		tokenHelper: t,
	}
}

//create an user
func (s *Service) CreateUser(name, phone, password string) (string, error) {
	//1. create new user entity
	newUserEntity, err := entity.NewUser(name, phone, password)
	if err != nil {
		return "", err
	}

	//2. search user with that phone, if there is exist, then return error
	found, err := s.userRepo.CheckUserExists(phone)
	if err != nil {
		return "", err
	}
	if found {
		return "", failure.ErrDuplicatePhone()
	}

	//3. call user repo to create new user on database
	return s.userRepo.Create(newUserEntity)
}

//authenticate an user
func (s *Service) Authenticate(phone, password string) (*entity.User, error) {
	//1. search user that have client request phone number
	userFound, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, failure.ErrUserNotFound() //whatever error in userrepo, it should return user not found
	}

	//2. validate entity
	if err := userFound.Validate(); err != nil {
		return nil, err
	}

	//3. verify password
	if err := userFound.VerifyPassword(password); err != nil {
		return nil, err
	}

	//4. generate access token & refresh token
	accessToken, refreshToken, err := s.tokenHelper.GenerateAllTokens(userFound.Name, userFound.Phone, userFound.ID)
	if err != nil {
		return nil, err
	}

	//5. set access token & refresh token
	userFound.SetToken(accessToken, refreshToken)

	return userFound, nil
}

//used to refresh user's token
func (s *Service) RefreshToken(refreshToken string) (*entity.User, error) {

	//1. check refresh token is still valid or not
	claims, err := s.tokenHelper.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	//2. get userObjectId from claims
	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return nil, failure.ErrConvertObjectIdToHex()
	}

	//3. get user information from database according to userObjectId
	userFound, err := s.userRepo.FindById(userObjectId)
	if err != nil {
		return nil, err //whatever error in userrepo, it should return user not found
	}

	//4. validate entity
	if err := userFound.Validate(); err != nil {
		return nil, err
	}

	//5. generate access token & refresh token
	accessToken, refreshToken, err := s.tokenHelper.GenerateAllTokens(userFound.Name, userFound.Phone, userFound.ID)
	if err != nil {
		return nil, err
	}

	//6. set access token & refresh token
	userFound.SetToken(accessToken, refreshToken)

	return userFound, nil
}

//used to check user's token
func (s *Service) CheckToken(token string) (*entity.User, error) {

	//1. check refresh token is still valid or not
	claims, err := s.tokenHelper.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	//2. get userObjectId from claims
	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return nil, failure.ErrConvertObjectIdToHex()
	}

	//3. get user information from database according to userObjectId
	userFound, err := s.userRepo.FindById(userObjectId)
	if err != nil {
		return nil, failure.ErrUserNotFound() //whatever error in userrepo, it should return user not found
	}

	//4. validate entity
	if err := userFound.Validate(); err != nil {
		return nil, err
	}

	return userFound, nil
}
