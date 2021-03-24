package hub

import (
	"encoding/json"
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
			_, msg, err := c.ws.ReadMessage() //waiting message from websocket
			if err != nil {
				c.write(websocket.CloseMessage, []byte{})
				log.Printf("[close] %v", err)
				break
			}

			//convert plain message data to formated message struct
			//example messageMap :
			//{"data":"Hi","fromUserId":"605ae53dce933ec8b23f9cc1","toUserId":"605ae3f2ce933ec8b23f9cbd","contactId":"605ae6dcdbadf9c66aa4fe60"}
			var message models.Message

			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("Cannot unmarshall message : %s", msg)
				break
			}

			MainHub.Broadcast <- message

		}

	}

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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case messageInByte, ok := <-c.send: //if there is message goes to channel send in certain connection, then send that message to websocket
			if !ok {
				c.write(websocket.CloseMessage, []byte{}) //write to ws indicate that connection was closed
				return
			}
			err := c.write(websocket.TextMessage, messageInByte) //write to ws with message : messageInByte
			if err != nil {
				return
			}

		case <-ticker.C:
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

	//3. make subscription instance (this explain that we make subscription with user : userId and con : c)
	s := subscription{c, userID}

	//4. register this subscription to hub
	MainHub.Register <- s
	//5. run this function on background
	go s.writePump()
	go s.readPump()
}
