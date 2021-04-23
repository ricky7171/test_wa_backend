package websocket

import (
	"encoding/json"

	"github.com/ricky7171/test_wa_backend/internal/entity"

	ChatUseCase "github.com/ricky7171/test_wa_backend/internal/usecase/chat"
)

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
	Broadcast chan entity.Message

	// Register requests from the connections.
	Register chan subscription

	// Unregister requests from connections.
	Unregister chan subscription
}

var MainHub = Hub{
	Broadcast:  make(chan entity.Message),
	Register:   make(chan subscription),
	Unregister: make(chan subscription),
	Users:      make(map[string]*connection),
}

func (h *Hub) Run(service *ChatUseCase.Service) {
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
			messageSaved, err := service.CreateMessage(&m)

			//2. convert object m to byte[]
			dataSend, err := json.Marshal(messageSaved)
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
				h.Users[m.ToUserId].send <- dataSend
			}

		}
	}
}
