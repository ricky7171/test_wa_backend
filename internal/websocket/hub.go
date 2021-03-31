package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ricky7171/test_wa_backend/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func SaveChat(msg models.Message, dbInstance *mongo.Database) error {
	//1. make ctx object
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. make chat model
	var chat models.Chat
	chat.ContactId, _ = primitive.ObjectIDFromHex(msg.ContactId)
	chat.CreatedAt = time.Now()
	chat.ID = primitive.NewObjectID()
	chat.Message = msg.Data
	chat.SenderId, _ = primitive.ObjectIDFromHex(msg.FromUserId)

	//3. insert new chat
	_, err := dbInstance.Collection("chats").InsertOne(ctx, chat)
	defer cancel()

	if err != nil {
		fmt.Println(err)
	}
	defer cancel()

	return nil

}

func (h *Hub) Run(dbInstance *mongo.Database) {
	for {
		select {
		case s := <-h.Register: //when there is user connect to ws
			if connNow, ok := h.Users[s.userId]; ok { //if there is other connection used by this user, then send warning message to that connection
				connNow.send <- []byte("This connection is lost, because you have opened a chat on another page")
			}
			h.Users[s.userId] = s.conn //fill ws connection to this user
		case s := <-h.Unregister: //when there is user disconnect from ws
			if _, ok := h.Users[s.userId]; ok {
				close(s.conn.send)

				delete(h.Users, s.userId)
			}
		case m := <-h.Broadcast: //when there is message

			//1. store to database
			SaveChat(m, dbInstance)

			//2. convert object m ke byte[]
			dataSend, err := json.Marshal(m)
			if err != nil {
				panic(err)
			}

			//3. send to ws
			//3.a. send to sender ws
			if _, ok := h.Users[m.FromUserId]; ok {
				h.Users[m.FromUserId].send <- dataSend
			} else {
				break
			}

			//3.b. send to receiver ws
			if _, ok := h.Users[m.ToUserId]; ok {
				h.Users[m.ToUserId].send <- dataSend //send message to receiver through "send" channel
			}

		}
	}
}
