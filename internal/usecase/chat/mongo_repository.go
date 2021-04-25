package chat

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

//used to get chats by contact ID
func (r *MongoRepository) FindByContact(contactId primitive.ObjectID, lastId string) ([]entity.Chat, error) {
	//1. make ctx with timout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. get data chat according to that contact id and last id
	matchContactIDPipeline := bson.D{{"$match", bson.D{{"contactId", contactId}}}}
	sortPipeline := bson.D{{"$sort", bson.D{{"_id", -1}}}}

	limitPipeline := bson.D{{"$limit", 20}}
	groupPipeline := bson.D{{"$project", bson.D{{"_id", 1}, {"contactId", 1}, {"senderId", 1}, {"message", 1}, {"createdAt", 1}}}}

	var cursor *mongo.Cursor
	var err error
	if lastId == "nil" { //means get first page (first 20 chat data)
		cursor, err = r.db.Collection("chats").Aggregate(
			ctx,
			mongo.Pipeline{
				matchContactIDPipeline,
				sortPipeline,
				limitPipeline,
				groupPipeline,
			},
		)
	} else { //means get another page (ex : 21 - 40, 41 - 60,etc)
		lastObjectID, err := primitive.ObjectIDFromHex(lastId)
		if err != nil {
			return nil, failure.ErrConvertObjectIdToHex()
		}

		paginatePipeline := bson.D{{"$match", bson.D{{"_id", bson.D{{"$lt", lastObjectID}}}}}}
		cursor, err = r.db.Collection("chats").Aggregate(
			ctx,
			mongo.Pipeline{
				matchContactIDPipeline,
				sortPipeline,
				paginatePipeline,
				limitPipeline,
				groupPipeline,
			},
		)
	}
	defer cancel()

	if err != nil {
		return nil, failure.ErrRepoFailedQueryGet()
	}

	//4. convert result to array of chat entity
	var chats []entity.Chat
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, failure.ErrRepoFailedConvert()
	}
	return chats, nil

}

//used to insert new chat data
func (r *MongoRepository) Create(message *entity.Message) (*entity.Message, error) {
	//1. make ctx with timeout 100 second
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. make chat entity
	var chat entity.Chat
	chat.ContactId, _ = primitive.ObjectIDFromHex(message.ContactId)
	chat.CreatedAt = time.Now()
	chat.ID = primitive.NewObjectID()
	chat.Message = message.Data
	chat.SenderId, _ = primitive.ObjectIDFromHex(message.FromUserId)

	//3. validate entity
	if err := chat.Validate(); err != nil {
		return nil, err
	}

	//4. insert new chat
	_, err := r.db.Collection("chats").InsertOne(ctx, chat)
	defer cancel()

	if err != nil {
		return nil, failure.ErrRepoFailedQueryInsert()
	}

	return message, nil
}
