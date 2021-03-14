package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"wa/models"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	//timeout saat ngirim data
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// keknya sih ini timeout dalam room nya
	// nilai pingPeriod ini adalah 54 detik
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
	mu   sync.Mutex
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		MainHub.Unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("error saat setReadDeadline")
	} else {
		c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			_, msg, err := c.ws.ReadMessage() //nungguin message masuk dari websocket
			if err != nil {
				c.write(websocket.CloseMessage, []byte{})
				log.Printf("error: %v", err)
				break
			}

			//convert plain message data to formated message struct
			//example messageMap :
			//{"data" : "bla bla bla" in byte, "fromUserId" : 1, "toUserId" : 2}
			var messageMap map[string]interface{}
			if err := json.Unmarshal(msg, &messageMap); err != nil {
				panic(err)
			}

			//convert semua key di map messageMap. Kalau ada yg error, langsung break
			data, ok := messageMap["data"].(string)
			if !ok {
				break
			}
			fromUserID, ok := messageMap["fromUserId"].(string)
			if !ok {
				break
			}
			toUserID, ok := messageMap["toUserId"].(string)
			if !ok {
				break
			}
			roomID, ok := messageMap["roomId"].(string)
			if !ok {
				break
			}

			//build model message
			m := models.Message{Data: data, FromUserId: fromUserID, ToUserId: toUserID, Room_id: roomID}

			//kirim ke channel broadcast
			MainHub.Broadcast <- m

		}

	}

}

// write writes a message with the given message type and payload.
//mt itu messageType
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	fmt.Println("masuk writepump")
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case messageInByte, ok := <-c.send: //jika ada pesan masuk ke dalam channel send di koneksi tertentu, maka kirim pesan tsb ke websocket
			if !ok {
				c.write(websocket.CloseMessage, []byte{}) //write ke ws, tipe write nya adalah close message, isinya adalah {}
				return
			}
			err := c.write(websocket.TextMessage, messageInByte)
			if err != nil { //coba write ke ws, tipenya adalah textmessage, isinya adalah variabel message, kalau tidak error maka
				return
			}

		case <-ticker.C: //keknya, kalau roomnya timeout maka
			err := c.write(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}

}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, userID string) {

	//1. upgrade protocol from http to websocket
	//then, save the connection to variable ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//2. just make connection instance and save to variable c
	c := &connection{send: make(chan []byte, 256), ws: ws}

	//3. make subscription instance (this explain that we make subscription with room : roomId and con : c)
	s := subscription{c, userID}

	//4. register this subscription to hub
	MainHub.Register <- s

	//5. run this function on background
	go s.writePump()
	go s.readPump()
	fmt.Println("tidak ada masalah di serve")
}
