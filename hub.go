package main

//contoh message
// {data : "blablablabla", room : "room_a"}
// ini masih belum tau buat apa, tapi intinya ini adalah data antara message dengan room
type message struct {
	data []byte
	room string
}

//contoh subscription
// {con : con1, room : "room_a"}
// jadi semacam data antara room dan koneksi
// baik register maupun unregister, tipe datanya ya subscription
type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

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
	rooms:      make(map[string]map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register: //kalau ada subs yg register
			connections := h.rooms[s.room] //cek koneksi di room tersebut
			if connections == nil {        //kalau belum ada koneksi atau room tsb masih belum ada, maka inisialisasi koneksi KOSONG baru di room tsb
				connections = make(map[*connection]bool) //inisialisasi koneksi baru
				h.rooms[s.room] = connections            //masukan koneksi baru tersebut ke room tadi
			}
			h.rooms[s.room][s.conn] = true //catat, di room tersebut dengan koneksi tsb (yg didapat dari subs) sudah hidup / true.
		case s := <-h.unregister: //kalau ada subs yang unregister
			connections := h.rooms[s.room] //cek koneksi di room tersebut
			if connections != nil {        //kalau sudah ada koneksinya
				if _, ok := connections[s.conn]; ok { //kalau koneksinya hidup / true
					delete(connections, s.conn) //delete key s.conn di variabel connections
					close(s.conn.send)          //close channel yang ada di s.conn
					if len(connections) == 0 {  //kalau koneksi di room tersebut sudah habis alias 0
						delete(h.rooms, s.room) //hapus aja room tersebut dari list rooms yang ada di h.rooms
					}
				}
			}
		case m := <-h.broadcast: //kalau ada broadcast/message masuk
			connections := h.rooms[m.room]
			for c := range connections { //loop semua koneksi yang ada di room tersebut
				select {
				case c.send <- m.data: //kirim pesan ke semua koneksi yang ada di room tersebut (melalui channel send)
				default: //kalau gagal (mungkin udah error koneksinya, mending hapus aja)
					close(c.send)              //close channel send
					delete(connections, c)     //delete key c di variabel connections
					if len(connections) == 0 { //kalau koneksi di room tersebut sudah habis alias 0
						delete(h.rooms, m.room) //hapus aja room tersebut dari list rooms yang ada di h.rooms
					}
				}
			}
		}
	}
}
