package websocket

import (
	"log"
)

type Hub struct {
	RoomId string
	clients map[*client]bool
	broadcast chan []byte
	register chan *client
	unregister chan *client
}

func NewHub(roomId string) *Hub {
	return &Hub{
		RoomId:     roomId,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected to room '%s'. Room size: %d", h.RoomID, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)   
				log.Printf("Client disconnected from room '%s'. Room size: %d", h.RoomID, len(h.clients))
			}
			
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}