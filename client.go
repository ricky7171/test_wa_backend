package main

import (
	"log"
	"net/http"
	"sync"
	"time"

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
		//fmt.Println("defer di readpump")
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		//log.Printf("error saat setReadDeadline")
	} else {
		c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			c.mu.Lock()
			_, msg, err := c.ws.ReadMessage()
			c.mu.Unlock()
			if err != nil {
				//fmt.Println("exit")
				//c.ws.Close()
				c.mu.Lock()
				c.write(websocket.CloseMessage, []byte{})
				c.mu.Unlock()
				//log.Printf("ada error nih waktu nge read message dari ws")
				//log.Printf("error: %v", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					//log.Printf("error waktu nge read : unexpectedcloseerror")
					//log.Printf("error: %v", err)
				} else {
					// log.Printf("ada error nih waktu nge read message dari ws")
					// log.Printf("error: %v", err)
				}
				break
			}
			m := message{msg, s.room}
			h.broadcast <- m
		}

	}

}

// write writes a message with the given message type and payload.
//mt itu messageType
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//fmt.Println("defer di writepump")
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send: //jika ada pesan masuk ke dalam channel send di koneksi tertentu, maka kirim pesan tsb ke websocket
			//fmt.Println("ada pesan masuk ke dalam room")
			//fmt.Println("ok atau tidak", ok)
			//fmt.Println("messagenya adalah : ", message)

			if !ok {
				c.mu.Lock()
				c.write(websocket.CloseMessage, []byte{}) //write ke ws, tipe write nya adalah close message, isinya adalah {}
				c.mu.Unlock()
				return
			}
			c.mu.Lock()
			err := c.write(websocket.TextMessage, message)
			c.mu.Unlock()
			if err != nil { //coba write ke ws, tipenya adalah textmessage, isinya adalah variabel message, kalau tidak error maka
				return
			}

		case <-ticker.C: //keknya, kalau roomnya timeout maka
			c.mu.Lock()
			err := c.write(websocket.PingMessage, []byte{})
			c.mu.Unlock()
			if err != nil {
				return
			}
		}
	}

}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, roomId string) {
	// fmt.Println("masuk fungsi serveWs")

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
	s := subscription{c, roomId}

	//4. register this subscription to hub
	// fmt.Println("send subscriber to register channel")
	h.register <- s

	//5. run this function on background
	go s.writePump()
	go s.readPump()
}
