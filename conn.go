package main

import (
	"fmt"
	"github.com/pborman/uuid"
	"golang.org/x/net/websocket"
	"io"
	"os"
)

type Conn struct {
	*websocket.Conn
	id   string
	room *Room
}

func NewConn(conn *websocket.Conn, room *Room) *Conn {
	return &Conn{conn, uuid.New(), room}
}

func (c *Conn) Send(msg *Message) {
	if err := websocket.JSON.Send(c.Conn, msg); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (c *Conn) Run() {
	msg := new(Message)

	for {
		msg.Clear()

		switch err := websocket.JSON.Receive(c.Conn, msg); {
		case err == io.EOF:
			return
		case err != nil:
			fmt.Fprintln(os.Stderr, err)
		default:
			msg.From = c.id
			c.room.Send(msg)
		}
	}
}
