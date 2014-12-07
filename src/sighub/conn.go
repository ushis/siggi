package sighub

import (
	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type Conn struct {
	*websocket.Conn
	id  string
	hub chan *Message
	die chan string
}

func NewConn(conn *websocket.Conn, hub chan *Message, die chan string) *Conn {
	return &Conn{conn, uuid.New(), hub, die}
}

func (c *Conn) Send(msg *Message) error {
	return websocket.JSON.Send(c.Conn, msg)
}

func (c *Conn) Run() {
	for {
		msg := NewMessage()

		switch err := websocket.JSON.Receive(c.Conn, msg); {
		case err == io.EOF:
			c.die <- c.id
			return
		case err != nil:
			log.Println(err)
		default:
			msg.From = c.id
			c.hub <- msg
		}
	}
}
