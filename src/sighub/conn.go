package sighub

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"os"
)

type Conn struct {
	*websocket.Conn
	id  string
	hub chan *Message
}

func NewConn(conn *websocket.Conn, hub chan *Message) *Conn {
	return &Conn{conn, uuid.New(), hub}
}

func (c *Conn) Send(msg *Message) error {
	return websocket.JSON.Send(c.Conn, msg)
}

func (c *Conn) Run() {
	for {
		msg := NewMessage()

		switch err := websocket.JSON.Receive(c.Conn, msg); {
		case err == io.EOF:
			return
		case err != nil:
			fmt.Fprintln(os.Stderr, err)
		default:
			msg.From = c.id
			c.hub <- msg
		}
	}
}
