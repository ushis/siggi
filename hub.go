package main

import (
	"golang.org/x/net/websocket"
)

type Hub struct {
	rooms map[string]Room
	msg   chan *Message
	reg   chan *Conn
	rm    chan *Conn
	die   chan int
}

func NewHub() *Hub {
	return &Hub{
		make(map[string]Room),
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
	c := NewConn(conn, h.msg)
	h.reg <- c
	c.Run()
	h.rm <- c
}

func (h *Hub) register(conn *Conn) {
	room, ok := h.rooms[conn.roomId]

	if !ok {
		room = NewRoom()
		h.rooms[conn.roomId] = room
	}
	room.Add(conn)
}

func (h *Hub) recv(msg *Message) {
	if room, ok := h.rooms[msg.Room]; ok {
		room.Send(msg)
	}
	msg.Free()
}

func (h *Hub) closeConn(conn *Conn) {
	conn.Close()

	if room, ok := h.rooms[conn.roomId]; ok {
		room.Rm(conn)

		if room.Len() == 0 {
			delete(h.rooms, conn.roomId)
		}
	}
}

func (h *Hub) shutdown() {
	for _, room := range h.rooms {
		room.Each(h.closeConn)
	}
	close(h.msg)
	close(h.reg)
	close(h.rm)
	close(h.die)
}
