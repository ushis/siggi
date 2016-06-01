package main

import (
	"sync"
)

type Room struct {
	conns map[string]*Conn
	mutex *sync.RWMutex
}

func NewRoom() *Room {
	return &Room{conns: make(map[string]*Conn), mutex: &sync.RWMutex{}}
}

func (r *Room) Add(conn *Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.conns[conn.id] = conn
}

func (r *Room) Rm(conn *Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.conns, conn.id)
	conn.Close()
}

func (r *Room) Get(id string) (*Conn, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if conn, ok := r.conns[id]; ok {
		return conn, true
	}
	return nil, false
}

func (r *Room) Len() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.conns)
}

func (r *Room) IsEmpty() bool {
	return r.Len() == 0
}

func (r *Room) Close() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, conn := range r.conns {
		delete(r.conns, conn.id)
		conn.Close()
	}
}

func (r *Room) Send(msg *Message) {
	if len(msg.To) == 0 {
		r.broadcast(msg)
	} else {
		r.send(msg)
	}
}

func (r *Room) broadcast(msg *Message) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, conn := range r.conns {
		if conn.id != msg.From {
			conn.Send(msg)
		}
	}
}

func (r *Room) send(msg *Message) {
	if conn, ok := r.Get(msg.To); ok {
		conn.Send(msg)
	}
}
