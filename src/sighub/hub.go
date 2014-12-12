package sighub

import (
	"golang.org/x/net/websocket"
)

type Hub struct {
	rooms map[string]map[string]*Conn
	msg   chan *Message
	reg   chan *Conn
	rm    chan *Conn
	die   chan int
}

func New() *Hub {
	return &Hub{
		make(map[string]map[string]*Conn),
		make(chan *Message),
		make(chan *Conn),
		make(chan *Conn),
		make(chan int),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.msg:
			h.recv(msg)
		case conn := <-h.reg:
			h.register(conn)
		case conn := <-h.rm:
			h.closeConn(conn)
		case <-h.die:
			h.shutdown()
			return
		}
	}
}

func (h *Hub) Die() {
	h.die <- 1
}

func (h *Hub) HTTPHandler() websocket.Handler {
	return websocket.Handler(h.connect)
}

func (h *Hub) connect(conn *websocket.Conn) {
	c := NewConn(conn, conn.Config().Location.Query().Get("room"), h.msg)
	h.reg <- c
	c.Run()
	h.rm <- c
}

func (h *Hub) register(conn *Conn) {
	room, ok := h.rooms[conn.room]

	if !ok {
		room = make(map[string]*Conn)
		h.rooms[conn.room] = room
	}
	room[conn.id] = conn
}

func (h *Hub) recv(msg *Message) {
	if len(msg.To) > 0 {
		h.send(msg)
	} else {
		h.broadcast(msg)
	}
	msg.Free()
}

func (h *Hub) broadcast(msg *Message) {
	if room, ok := h.rooms[msg.Room]; ok {
		for id, conn := range room {
			if id != msg.From {
				conn.Send(msg)
			}
		}
	}
}

func (h *Hub) send(msg *Message) {
	if room, ok := h.rooms[msg.Room]; ok {
		if conn, ok := room[msg.To]; ok {
			conn.Send(msg)
		}
	}
}

func (h *Hub) closeConn(conn *Conn) {
	conn.Close()

	room, ok := h.rooms[conn.room]

	if !ok {
		return
	}
	delete(room, conn.id)

	if len(room) == 0 {
		delete(h.rooms, conn.room)
	}
}

func (h *Hub) shutdown() {
	for _, room := range h.rooms {
		for _, conn := range room {
			h.closeConn(conn)
		}
	}
	close(h.msg)
	close(h.reg)
	close(h.rm)
	close(h.die)
}
