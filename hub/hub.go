package hub

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"wa/database"
	"wa/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var roomCollection *mongo.Collection = database.OpenCollection(database.Client, "rooms")
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

//subscription example
// {con : con1, userId : xxx123}
type subscription struct {
	conn   *connection
	userId string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	//key is id and value is connection ws
	Users map[string]*connection

	// Inbound messages from the connections.
	Broadcast chan models.Message

	// Register requests from the connections.
	Register chan subscription

	// Unregister requests from connections.
	Unregister chan subscription
}

var MainHub = Hub{
	Broadcast:  make(chan models.Message),
	Register:   make(chan subscription),
	Unregister: make(chan subscription),
	Users:      make(map[string]*connection),
}

func SaveChat(msg models.Message) error {
	//1. make ctx object
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. fill created_at
	msg.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//3. if room_id not found in msg object, then find room that used by sender and receiver
	//if room not found, then make new room
	newFromUserIdObject, err := primitive.ObjectIDFromHex(msg.FromUserId)
	newToUserIdObject, err := primitive.ObjectIDFromHex(msg.ToUserId)

	makeNewRoom := false //=> means, no need to make new room
	if msg.Room_id == "" {
		var roomExists models.Room
		opts := options.FindOne().SetProjection(bson.M{
			"room_id": 1,
		})
		err := roomCollection.FindOne(
			ctx,
			bson.M{
				"$or": []interface{}{
					bson.M{
						"$and": []interface{}{
							bson.M{
								"users.0": newFromUserIdObject,
							},
							bson.M{
								"users.1": newToUserIdObject,
							},
						},
					},
					bson.M{
						"$and": []interface{}{
							bson.M{
								"users.1": newFromUserIdObject,
							},
							bson.M{
								"users.0": newToUserIdObject,
							},
						},
					},
				},
			},
			opts,
		).Decode(&roomExists)
		defer cancel()
		if err != nil { //means, room not found, then need to make new room
			makeNewRoom = true
		} else { //means, room found, then NO need to make new room
			msg.Room_id = roomExists.ID.Hex()
		}

	}

	//4. if room found, then update field chat_history and insert chat in that field
	if !makeNewRoom {
		roomObjectID, _ := primitive.ObjectIDFromHex(msg.Room_id)
		_, err := roomCollection.UpdateOne(ctx, bson.M{
			"_id": roomObjectID,
		},
			bson.M{
				"$push": bson.M{
					"chat_history": bson.M{
						"_id":        primitive.NewObjectID(),
						"user_id":    msg.FromUserId,
						"message":    msg.Data,
						"created_at": msg.Created_at,
					},
				},
			})
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
	} else { //5. if room not found, then insert new room also with that chat inside
		newRoomId := primitive.NewObjectID()

		if err != nil {
			return nil
		}

		if err != nil {
			return nil
		}
		roomCollection.InsertOne(ctx, bson.M{
			"_id": newRoomId,
			"users": []interface{}{
				newFromUserIdObject,
				newToUserIdObject,
			},
			"chat_history": []interface{}{
				bson.M{
					"_id":        primitive.NewObjectID(),
					"user_id":    msg.FromUserId,
					"message":    msg.Data,
					"created_at": msg.Created_at,
				},
			},
		})
		defer cancel()
	}

	return nil

}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.Register: //when there is user connect to ws
			h.Users[s.userId] = s.conn //fill h.users with key that user id, then fill value with ws connection
		case s := <-h.Unregister: //when there is user disconnect from ws
			if _, ok := h.Users[s.userId]; ok {
				close(s.conn.send)
				delete(h.Users, s.userId)
			}
		case m := <-h.Broadcast: //when there is message

			//1. store to database
			go SaveChat(m)

			//2. convert object m ke byte[]
			dataSend, err := json.Marshal(m)
			if err != nil {
				panic(err)
			}

			//3. send to ws
			//3.a. send to sender ws
			if _, ok := h.Users[m.FromUserId]; ok {
				select {
				case h.Users[m.FromUserId].send <- dataSend: //send message to sender through "send" channel
				default: //this is when fail to send the message (maybe user exit from browser)
					break
				}
			} else {
				break
			}

			//3.b. send to receiver ws
			if _, ok := h.Users[m.ToUserId]; ok {
				select {
				case h.Users[m.ToUserId].send <- dataSend: //send message to receiver through "send" channel
				default: //this is when fail to send the message (maybe user exit from browser)
				}
			} else {
			}

		}
	}
}
