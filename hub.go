package main

import (
	"encoding/json"
)

//contoh message
// {data : "blablablabla", fromUserId : 1, toUserId : 2}
// ini masih belum tau buat apa, tapi intinya ini adalah data antara message dengan room
type message struct {
	data       []byte
	fromUserId int
	toUserId   int
}

//contoh subscription
// {con : con1, userId : 1}
// jadi semacam data antara room dan koneksi
// baik register maupun unregister, tipe datanya ya subscription
type subscription struct {
	conn   *connection
	userId int
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	//key is id and value is connection ws
	users map[int]*connection

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	users:      make(map[int]*connection),
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register: //kalau ada subs yg register (dg kata lain, kalau ada user baru yg konek ke ws)
			h.users[s.userId] = s.conn //kalaupun sebelumnya sudah terisi, ya langsung ke replace
		case s := <-h.unregister: //kalau ada subs yang unregister (dg kata lain, kalau ada user yang exit dari koneksi ws)
			if _, ok := h.users[s.userId]; ok {
				close(s.conn.send)
				delete(h.users, s.userId)
			}
		case m := <-h.broadcast: //kalau ada broadcast/message masuk
			dataSend, err := json.Marshal(m)
			if err != nil {
				panic(err)
			}

			select {
			case h.users[m.toUserId].send <- dataSend: //kirim pesan ke penerimanya (melalui channel send)
				//store to database
			default: //kalau gagal (mungkin penerimanya error koneksinya, mending hapus aja)
				close(h.users[m.toUserId].send) //close channel send
				delete(h.users, m.toUserId)
			}

		}
	}
}
