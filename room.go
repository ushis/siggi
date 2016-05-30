package main

type Room map[string]*Conn

func NewRoom() Room {
	return make(Room)
}

func (r Room) Add(conn *Conn) {
	r[conn.id] = conn
}

func (r Room) Rm(conn *Conn) {
	delete(r, conn.id)
}

func (r Room) Len() int {
	return len(r)
}

func (r Room) Each(fn func(*Conn)) {
	for _, conn := range r {
		fn(conn)
	}
}

func (r Room) Send(msg *Message) {
	if len(msg.To) > 0 {
		r.send(msg)
	} else {
		r.broadcast(msg)
	}
}

func (r Room) Close() {
	for _, conn := range r {
		conn.Close()
		r.Rm(conn)
	}
}

func (r Room) send(msg *Message) {
	if conn, ok := r[msg.To]; ok {
		conn.Send(msg)
	}
}

func (r Room) broadcast(msg *Message) {
	for _, conn := range r {
		if conn.id != msg.From {
			conn.Send(msg)
		}
	}
}
