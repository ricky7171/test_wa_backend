package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DBinstance func
func DBinstance() *mongo.Database {

	//1. read mongoDB url
	MongoDb := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal("Error open db : ", err)
	}

	//2. make context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//3. connect to DB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	//4. make dbInstance and return it
	dbInstance := client.Database(os.Getenv("DATABASE_NAME"))
	return dbInstance
}

//OpenCollection is a function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	nameDb := os.Getenv("DATABASE_NAME")
	var collection *mongo.Collection = client.Database(nameDb).Collection(collectionName)

	return collection
}
