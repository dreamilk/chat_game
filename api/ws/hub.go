package ws

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	m map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{}
}

func Worker() {
}
