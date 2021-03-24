package hub

import (
	"context"
	"encoding/json"
	"log"
	"time"
	db "wa/database"
	"wa/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func SaveChat(msg models.Message) error {
	//1. make ctx object
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. fill created_at
	msg.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//3. if contact_id not found in msg object, then find contact first
	//if contact not found, then make new contact
	newFromUserIdObject, err := primitive.ObjectIDFromHex(msg.FromUserId)
	newToUserIdObject, err := primitive.ObjectIDFromHex(msg.ToUserId)

	makeNewContact := false //=> means, no need to make new contact
	if msg.Contact_id == "" {
		var contactExists models.Contact
		err := db.ContactCollection.FindOne(
			ctx,
			bson.M{
				"users": bson.M{
					"$all": []interface{}{
						newFromUserIdObject,
						newToUserIdObject,
					},
				},
			},
		).Decode(&contactExists)
		defer cancel()
		if err != nil { //means, contact not found, then need to make new contact
			makeNewContact = true
		} else { //means, contact found, then NO need to make new contact
			msg.Contact_id = contactExists.ID.Hex()
		}

	}

	//4. if contact found, then set contact_id according to msg.Contact_id
	var chat models.Chat
	if !makeNewContact {
		chat.Contact_id, _ = primitive.ObjectIDFromHex(msg.Contact_id)
	} else { //5. if contact not found, then insert new contact
		newContactId := primitive.NewObjectID()
		chat.Contact_id = newContactId
		_, err := db.ContactCollection.InsertOne(ctx, bson.M{
			"_id": newContactId,
			"users": []interface{}{
				newFromUserIdObject,
				newToUserIdObject,
			},
		})
		defer cancel()

		if err != nil {
			log.Fatal(err)
			return nil
		}
		defer cancel()
	}

	//6. insert new chat
	chat.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	chat.ID = primitive.NewObjectID()
	chat.Message = msg.Data
	chat.Sender_id, _ = primitive.ObjectIDFromHex(msg.FromUserId)

	_, err = db.ChatCollection.InsertOne(ctx, chat)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

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
			SaveChat(m)

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
