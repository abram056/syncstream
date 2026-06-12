package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/abram056/syncstream/backend/internal/room"
)

// Hub manages WebSocket connections for a room.
type Hub struct {
	roomID     string
	manager    *room.Manager
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new hub for a room.
func NewHub(roomID string, manager *room.Manager) *Hub {
	return &Hub{
		roomID:     roomID,
		manager:    manager,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("client registered in room %s, total: %d", h.roomID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				h.mu.Unlock()

				// remove participant from room
				if client.ParticipantID != "" {
					if err := h.manager.RemoveParticipant(h.roomID, client.ParticipantID); err != nil {
						log.Printf("failed to remove participant: %v", err)
					}

					// broadcast user_left to remaining clients
					evt := map[string]interface{}{
						"type":    "user_left",
						"user_id": client.ParticipantID,
					}
					if msg, err := json.Marshal(evt); err == nil {
						h.Broadcast(msg)
					}
				}

				log.Printf("client unregistered from room %s, total: %d", h.roomID, len(h.clients))
			} else {
				h.mu.Unlock()
			}

		case msg := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- msg:
				default:
					// client send channel full, close it
					go func(c *Client) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients in the room.
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// RegisterClient registers a client with the hub.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client from the hub.
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// HandleJoinRoom processes a join_room event.
func (h *Hub) HandleJoinRoom(client *Client, evt map[string]interface{}) error {
	displayName, ok := evt["display_name"].(string)
	if !ok || displayName == "" {
		return ErrInvalidEvent
	}

	// add participant to room via manager
	participant, err := h.manager.AddParticipant(h.roomID, client.ParticipantID, displayName)
	if err != nil {
		return err
	}

	client.DisplayName = participant.DisplayName

	// send room_joined confirmation to the client
	roomJoined := map[string]interface{}{
		"type":    "room_joined",
		"room_id": h.roomID,
	}
	if msg, err := json.Marshal(roomJoined); err == nil {
		client.Send <- msg
	}

	// send room_state to the client
	r, err := h.manager.GetRoomByID(h.roomID)
	if err != nil {
		return err
	}

	roomState := map[string]interface{}{
		"type":         "room_state",
		"room_id":      r.ID,
		"status":       string(r.Status),
		"media_url":    r.Media.URL,
		"is_playing":   r.PlaybackState.IsPlaying,
		"position":     r.PlaybackState.Position,
		"participants": len(r.Participants),
	}
	if msg, err := json.Marshal(roomState); err == nil {
		client.Send <- msg
	}

	// broadcast user_joined to other clients
	userJoined := map[string]interface{}{
		"type":         "user_joined",
		"user_id":      client.ParticipantID,
		"display_name": displayName,
	}
	if msg, err := json.Marshal(userJoined); err == nil {
		h.Broadcast(msg)
	}

	return nil
}
