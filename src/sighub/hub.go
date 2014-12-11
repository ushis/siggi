package sighub

import (
	"fmt"
	"golang.org/x/net/websocket"
	"os"
)

type Hub struct {
	conns map[string]*Conn
	msg   chan *Message
	reg   chan *Conn
	rm    chan *Conn
	die   chan int
}

func New() *Hub {
	return &Hub{
		make(map[string]*Conn),
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
			h.conns[conn.id] = conn
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

func (h *Hub) recv(msg *Message) {
	if len(msg.To) > 0 {
		h.sendTo(msg.To, msg)
	} else {
		h.broadcast(msg)
	}
	msg.Free()
}

func (h *Hub) broadcast(msg *Message) {
	for id, conn := range h.conns {
		if msg.Room == conn.room && id != msg.From {
			h.sendTo(id, msg)
		}
	}
}

func (h *Hub) sendTo(id string, msg *Message) {
	if conn, ok := h.conns[id]; ok && conn.room == msg.Room {
		if err := conn.Send(msg); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (h *Hub) closeConn(conn *Conn) {
	conn.Close()
	delete(h.conns, conn.id)
}

func (h *Hub) shutdown() {
	for _, conn := range h.conns {
		h.closeConn(conn)
	}
	close(h.msg)
	close(h.reg)
	close(h.rm)
	close(h.die)
}
