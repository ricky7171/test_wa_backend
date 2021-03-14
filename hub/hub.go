package hub

import (
	"context"
	"encoding/json"
	"fmt"
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

//contoh subscription
// {con : con1, userId : 1}
// jadi semacam data antara room dan koneksi
// baik register maupun unregister, tipe datanya ya subscription
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
	//1. buat ctx dengan timeout 100 detik
	//dimana ctx ini bakal digunakan saat nge query
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	//2. isi created_at nya dulu
	msg.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//3. jika didalam object msg belum ada room_id nya,
	//maka cek apakah antara id pengirim dan penerima itu sudah memiliki room sebelumnya
	//kalau belum ya berarti buat room baru
	makeNewRoom := false //=> artinya tidak usah bikin room baru
	if msg.Room_id == "" {
		var roomExists models.Room
		opts := options.FindOne().SetProjection(bson.M{
			"room_id": 1,
		})
		fmt.Println("cari room apakah sudah ada atau belum")
		err := roomCollection.FindOne(
			ctx,
			bson.M{
				"$or": []interface{}{
					bson.M{
						"$and": []interface{}{
							bson.M{
								"users.0.user_id": msg.FromUserId,
							},
							bson.M{
								"users.1.user_id": msg.ToUserId,
							},
						},
					},
					bson.M{
						"$and": []interface{}{
							bson.M{
								"users.1.user_id": msg.FromUserId,
							},
							bson.M{
								"users.0.user_id": msg.ToUserId,
							},
						},
					},
				},
			},
			opts,
		).Decode(&roomExists)
		defer cancel()

		if err != nil { //artinya roomnya kosong
			makeNewRoom = true
		} else { //kalau ternyata ketemu roomnya
			msg.Room_id = roomExists.Room_id
		}

	}

	//jika roomnya sudah ada maka, cukup update field chat_history, lalu insert data chat disitu
	if !makeNewRoom {
		roomCollection.UpdateOne(ctx, bson.M{
			"room_id": msg.Room_id,
		},
			bson.M{
				"$push": bson.M{
					"chat_history": bson.M{
						"user_id":    msg.FromUserId,
						"message":    msg.Data,
						"created_at": msg.Created_at,
					},
				},
			})
		defer cancel()
	} else { //jika roomnya belum ada, maka insert room baru dengan chat_history didalamnya
		newRoomId := primitive.NewObjectID()
		newFromUserIdObject, err := primitive.ObjectIDFromHex(msg.FromUserId)
		if err != nil {
			return nil
		}
		newToUserIdObject, err := primitive.ObjectIDFromHex(msg.ToUserId)
		if err != nil {
			return nil
		}
		roomCollection.InsertOne(ctx, bson.M{
			"_id":     newRoomId,
			"room_id": newRoomId.Hex(),
			"users": []interface{}{
				newFromUserIdObject,
				newToUserIdObject,
			},
			"chat_history": []interface{}{
				bson.M{
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
		case s := <-h.Register: //kalau ada subs yg register (dg kata lain, kalau ada user baru yg konek ke ws)
			h.Users[s.userId] = s.conn //kalaupun sebelumnya sudah terisi, ya langsung ke replace
		case s := <-h.Unregister: //kalau ada subs yang unregister (dg kata lain, kalau ada user yang exit dari koneksi ws)
			if _, ok := h.Users[s.userId]; ok {
				close(s.conn.send)
				delete(h.Users, s.userId)
			}
		case m := <-h.Broadcast: //kalau ada broadcast/message masuk

			//1. store to database dulu
			go SaveChat(m)

			//2. convert object m ke byte[]
			dataSend, err := json.Marshal(m)
			if err != nil {
				panic(err)
			}

			//kirim ke websocketnya
			//3. kirim ke ws nya pengirim dulu
			if _, ok := h.Users[m.FromUserId]; ok {
				select {
				case h.Users[m.FromUserId].send <- dataSend: //kirim pesan ke pengirim (melalui channel send)
				default: //kalau gagal (mungkin pengirimnya abis kirim pesan, langsung diexit browsernya)
					break
				}
			} else { //mungkin pengirimnya sudah tidak aktif lagi
				break
			}

			//4. baru kirim ke wsnya penerima
			if _, ok := h.Users[m.ToUserId]; ok {
				select {
				case h.Users[m.ToUserId].send <- dataSend: //kirim pesan ke penerimanya (melalui channel send)
				default: //kalau gagal (mungkin penerimanya error koneksinya, mending hapus aja)
				}
			} else { //mungkin penerimanya sudah tidak aktif lagi
			}

		}
	}
}
