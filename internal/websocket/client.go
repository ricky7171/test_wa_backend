package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ricky7171/test_wa_backend/internal/models"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
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
	c.ws.SetReadLimit(maxMessageSize)

	var duplicateDetected bool
	for {
		_, msg, err := c.ws.ReadMessage() //waiting message from websocket
		if err != nil {
			if !strings.Contains(err.Error(), "Duplicate connection") { //when there is an error, then check if the error is NOT because duplicate connection, then give flag that this is not duplicate
				duplicateDetected = false
				fmt.Println("[close] ", err.Error())
			} else { //if there is an error and this error because duplicate connection, then give flag that there is duplicate connection occur
				duplicateDetected = true
			}
			break //whatever the error, it will go out from this for loop

		}

		//convert plain message data to formated message struct
		//example messageMap :
		//{"data":"Hi","from_user_id":"605ae53dce933ec8b23f9cc1","to_user_id":"605ae3f2ce933ec8b23f9cbd","contact_id":"605ae6dcdbadf9c66aa4fe60"}
		var message models.Message

		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println("Cannot unmarshall message : ", msg)
			duplicateDetected = false //this means, that this error not because duplicate connection
			break
		}
		MainHub.Broadcast <- message
	}

	//notify client that websocket was closed
	c.write(websocket.CloseMessage, []byte{})

	if duplicateDetected { //if duplicate connection detected, then close only channel in this connection
		close(s.conn.send)
	} else { //if this is not because duplicate connection (may be user exit from browser), then unregister this subscription (means remove this subs from memory)
		MainHub.Unregister <- s
	}

	//close this websocket connection
	c.ws.Close()

}

// write writes a message with the given message type and payload.
//mt is messageType
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	defer func() {
		c.ws.Close()
	}()
	for {
		messageInByte, ok := <-c.send //block until there is message from channel "send" in certain connection, then send that message to websocket
		if !ok {
			c.write(websocket.CloseMessage, []byte{}) //write to ws indicate that connection was closed
			return
		}
		err := c.write(websocket.TextMessage, messageInByte) //write to ws with message : messageInByte
		if err != nil {
			return
		}

	}

}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, userID string) {
	//1. upgrade protocol from http to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//2. just make connection instance and save to variable c
	c := &connection{send: make(chan []byte, 256), ws: ws}

	//3. make subscription instance (this explain that we make subscription with user : userId and con : c)
	s := subscription{c, userID}

	//4. register this subscription to hub
	MainHub.Register <- s

	//5. run this function on background
	go s.writePump()
	go s.readPump()
}
