package contact

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

func (c *MongoRepository) FindByUser(firstUser *entity.UserWithName, secondUser *entity.UserWithName) ([]entity.Contact, error) {
	//1. make ctx with timout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. get all contact that contains this user id
	var inQuery bson.M
	if secondUser != nil {
		inQuery = bson.M{
			"$in": []entity.UserWithName{*firstUser, *secondUser},
		}
	} else {
		inQuery = bson.M{
			"$in": []entity.UserWithName{*firstUser},
		}
	}
	cursor, err := c.db.Collection("contacts").Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{
				{
					"$match", bson.M{
						"users": inQuery,
					},
				},
			},
		},
	)
	defer cancel()

	//3. convert cursor to contactWithName entity
	var allContacts []entity.Contact
	if err = cursor.All(ctx, &allContacts); err != nil {
		return nil, failure.ErrRepoFailedQueryGet()
	}
	return allContacts, nil

}

func (c *MongoRepository) Create(contact *entity.Contact) (string, error) {
	//1. make ctx with timeout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. insert contact to database
	resultInsert, err := c.db.Collection("contacts").InsertOne(ctx, contact)
	defer cancel()
	if err != nil {
		return "", failure.ErrRepoFailedQueryInsert()
	}

	//3. get new contact ID
	resultId := resultInsert.InsertedID.(primitive.ObjectID)

	return resultId.Hex(), nil
}
