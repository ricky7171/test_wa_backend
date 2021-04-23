package user

import (
	"context"
	"time"

	"github.com/ricky7171/test_wa_backend/internal/entity"
	"github.com/ricky7171/test_wa_backend/internal/failure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db *mongo.Database
}

//NewMongoRepository : factory to create MongoRepository object
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}

//used to insert new user data
func (r *MongoRepository) Create(user *entity.User) (string, error) {
	//1. make ctx with timeout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. run query insert to users's collection
	resultInsert, err := r.db.Collection("users").InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		return "", failure.ErrRepoFailedQueryInsert()
	}

	//3. get new user ID
	resultId := resultInsert.InsertedID.(primitive.ObjectID)

	//4. return user ID
	return resultId.Hex(), nil
}

//used to check duplicate specific user by phone
func (r *MongoRepository) CheckUserExists(phone string) (bool, error) {
	//1. make ctx with timeout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. run query get to user count
	count, err := r.db.Collection("users").CountDocuments(ctx, bson.M{"phone": phone})
	defer cancel()
	if err != nil {
		return false, failure.ErrRepoFailedQueryGet()
	}

	//3. check if count is 0
	if count == 0 {
		return false, nil
	}

	//4. return user ID
	return true, nil

}

//used to get specific user by phone
func (r *MongoRepository) FindByPhone(phone string) (*entity.User, error) {
	//1. make ctx with timeout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. run query get user data
	var foundUser entity.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"phone": phone}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		return nil, failure.ErrRepoFailedQueryGet()
	}

	//3. return user struct
	return &foundUser, nil

}

//used to get specific user by id
func (r *MongoRepository) FindById(userId string) (*entity.User, error) {
	//1. make ctx with timeout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. run query get user data
	var foundUser entity.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": userId}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		return nil, failure.ErrRepoFailedQueryGet()
	}

	//3. return user struct
	return &foundUser, nil

}
