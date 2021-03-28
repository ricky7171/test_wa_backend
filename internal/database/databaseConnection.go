package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DBinstance func
func DBinstance() *mongo.Database {

	//init .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal("Error open db : ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	dbInstance := client.Database(os.Getenv("DATABASE_NAME"))

	return dbInstance
}

//OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	nameDb := os.Getenv("DATABASE_NAME")
	var collection *mongo.Collection = client.Database(nameDb).Collection(collectionName)

	return collection
}

//var UserCollection *mongo.Collection = OpenCollection(Client, "users")
//var ContactCollection *mongo.Collection = OpenCollection(Client, "contacts")
//var ChatCollection *mongo.Collection = OpenCollection(Client, "chats")
