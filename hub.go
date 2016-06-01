package main

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Hub struct {
	rooms map[string]*Room
	mutex *sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{make(map[string]*Room), &sync.RWMutex{}}
}

func (h *Hub) HTTPHandler() websocket.Handler {
	return websocket.Handler(h.connect)
}

func (h *Hub) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for id, room := range h.rooms {
		delete(h.rooms, id)
		room.Close()
	}
}

func (h *Hub) connect(ws *websocket.Conn) {
	roomId := ws.Request().URL.Query().Get("room")
	room := h.getOrCreateRoom(roomId)
	conn := NewConn(ws, room)
	room.Add(conn)
	conn.Run()
	room.Rm(conn)
	h.cleanupRoom(roomId)
}

func (h *Hub) getOrCreateRoom(id string) *Room {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if room, ok := h.rooms[id]; ok {
		return room
	}
	room := NewRoom()
	h.rooms[id] = room
	return room
}

func (h *Hub) cleanupRoom(id string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if room, ok := h.rooms[id]; ok && room.IsEmpty() {
		delete(h.rooms, id)
	}
}
