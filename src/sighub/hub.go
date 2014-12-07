package sighub

import (
	"golang.org/x/net/websocket"
	"log"
)

type Hub struct {
	conns map[string]*Conn
	msg   chan *Message
	reg   chan *Conn
	rm    chan string
	die   chan int
}

func New() *Hub {
	return &Hub{
		make(map[string]*Conn),
		make(chan *Message),
		make(chan *Conn),
		make(chan string),
		make(chan int),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.msg:
			h.recv(msg)
		case conn := <-h.reg:
			h.conns[conn.id] = conn
		case id := <-h.rm:
			h.closeConn(id)
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
	c := NewConn(conn, h.msg, h.rm)
	h.reg <- c
	c.Run()
}

func (h *Hub) recv(msg *Message) {
	if len(msg.To) > 0 {
		h.sendTo(msg.To, msg)
	} else {
		h.broadcast(msg)
	}
	msg.Free()
}

func (h *Hub) broadcast(msg *Message) {
	for id := range h.conns {
		if id != msg.From {
			h.sendTo(id, msg)
		}
	}
}

func (h *Hub) sendTo(id string, msg *Message) {
	conn, ok := h.conns[id]

	if !ok {
		log.Printf("Invalid recipient id: %s\n", id)
		return
	}

	if err := conn.Send(msg); err != nil {
		log.Println(err)
	}
}

func (h *Hub) closeConn(id string) {
	h.conns[id].Close()
	delete(h.conns, id)
}

func (h *Hub) shutdown() {
	for id := range h.conns {
		h.closeConn(id)
	}
	close(h.msg)
	close(h.die)
}
