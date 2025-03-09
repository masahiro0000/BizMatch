package websocket

import (
	"log"
	"sync"
)

type BroadcastMessage struct {
	MatchID int64
	FromUser int64
	Content string
}

type Hub struct {
	clients map[int64]map[*Client]bool
	broadcast chan BroadcastMessage
	register chan *Client
	unregister chan *Client
	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]bool),
		broadcast: make(chan BroadcastMessage),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	// Continuously listen for events on the register, unregister, and broadcast channels.
	for {
		select {
		// Handle new client registration.
		case client := <-h.register:
			h.mu.Lock()
			// Initialize the clients map for the match if it doesn't already exist.
			if h.clients[client.MatchID] == nil {
				h.clients[client.MatchID] = make(map[*Client]bool)
			}
			h.clients[client.MatchID][client] = true
			h.mu.Unlock()
			log.Printf("Register client (user:%d) for matchID:%d\n", client.UserID, client.MatchID)

		// Handle client unregistration.
		case client := <-h.unregister:
			h.mu.Lock()
			// Check if the client's match exists in the clients map.
			if _, ok := h.clients[client.MatchID]; ok {
				delete(h.clients[client.MatchID], client)
				close(client.send)
				log.Printf("Unregister client (user:%d) from matchID:%d\n", client.UserID, client.MatchID)
				if len(h.clients[client.MatchID]) == 0 {
					delete(h.clients, client.MatchID)
				}
			}
			h.mu.Unlock()

		// Handle broadcasting messages to all clients in a match.
		case message := <-h.broadcast:
			h.mu.Lock()
			// Get the set of clients in the match for which the message is intended.
			clientsInRoom := h.clients[message.MatchID]
			for c := range clientsInRoom {
				select {
				// Message sent successfully.
				case c.send <-message:
				// If the send channel is blocked, close it and remove the client.
				default:
					close(c.send)
					delete(clientsInRoom, c)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends a broadcast message to the Hub's broadcast channel.
func (h *Hub) Broadcast(msg BroadcastMessage) {
	h.broadcast <-msg
}