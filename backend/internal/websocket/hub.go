package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

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

				// On websocket disconnect, mark participant as disconnected but keep them in the room.
				// They can reconnect within the grace period to resume their session.
				// If they do not reconnect in time, the participant expiration cleanup
				// will permanently remove them and broadcast user_left.
				if client.ParticipantID != "" {
					if err := h.manager.DisconnectParticipant(h.roomID, client.ParticipantID); err != nil {
						log.Printf("failed to disconnect participant: %v", err)
					}

					// broadcast user_disconnected to remaining clients
					evt := map[string]interface{}{
						"type":        "user_disconnected",
						"userId":      client.ParticipantID,
						"displayName": client.DisplayName,
					}
					if msg, err := json.Marshal(evt); err == nil {
						h.Broadcast(msg)
					}

					if err := h.broadcastRoomState("room_state"); err != nil {
						log.Printf("failed to broadcast room_state after disconnect: %v", err)
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
	displayName, ok := getDisplayName(evt)
	if !ok || displayName == "" {
		return ErrInvalidEvent
	}

	participant, err := h.manager.AddParticipant(h.roomID, client.ParticipantID, displayName)
	if err != nil {
		return err
	}

	client.DisplayName = participant.DisplayName

	roomJoined := map[string]interface{}{
		"type":   "room_joined",
		"roomId": h.roomID,
	}
	if msg, err := json.Marshal(roomJoined); err == nil {
		client.Send <- msg
	}

	r, err := h.manager.GetRoomByID(h.roomID)
	if err != nil {
		return err
	}

	roomState := map[string]interface{}{
		"type":              "room_state",
		"roomId":            r.ID,
		"status":            string(r.Status),
		"mediaUrl":          r.Media.URL,
		"isPlaying":         r.PlaybackState.IsPlaying,
		"position":          r.PlaybackState.Position,
		"numOfParticipants": len(r.Participants),
	}
	if msg, err := json.Marshal(roomState); err == nil {
		client.Send <- msg
	}

	userJoined := map[string]interface{}{
		"type":        "user_joined",
		"userId":      client.ParticipantID,
		"displayName": displayName,
	}
	if msg, err := json.Marshal(userJoined); err == nil {
		h.Broadcast(msg)
	}

	if err := h.broadcastRoomState("room_state"); err != nil {
		return err
	}

	return nil
}

func getDisplayName(evt map[string]interface{}) (string, bool) {
	if v, ok := evt["display_name"].(string); ok {
		return v, true
	}
	if v, ok := evt["displayName"].(string); ok {
		return v, true
	}
	return "", false
}

func (h *Hub) HandlePlay(client *Client, position float64, hasPosition bool) error {
	if hasPosition {
		if _, err := h.manager.Seek(h.roomID, client.ParticipantID, position); err != nil {
			return err
		}
	}

	if _, err := h.manager.Play(h.roomID, client.ParticipantID); err != nil {
		return err
	}

	return h.broadcastRoomState("sync_state")
}

func (h *Hub) HandlePause(client *Client, position float64, hasPosition bool) error {
	if hasPosition {
		if _, err := h.manager.Seek(h.roomID, client.ParticipantID, position); err != nil {
			return err
		}
	}

	if _, err := h.manager.Pause(h.roomID, client.ParticipantID); err != nil {
		return err
	}

	return h.broadcastRoomState("sync_state")
}

func (h *Hub) HandleSeek(client *Client, position float64) error {
	if _, err := h.manager.Seek(h.roomID, client.ParticipantID, position); err != nil {
		return err
	}

	return h.broadcastRoomState("sync_state")
}

func (h *Hub) broadcastRoomState(eventType string) error {
	r, err := h.manager.GetRoomByID(h.roomID)
	if err != nil {
		return err
	}

	roomState := map[string]interface{}{
		"type":              eventType,
		"roomId":            r.ID,
		"status":            string(r.Status),
		"mediaUrl":          r.Media.URL,
		"isPlaying":         r.PlaybackState.IsPlaying,
		"position":          r.PlaybackState.Position,
		"updatedBy":         r.PlaybackState.UpdatedBy,
		"updatedAt":         r.PlaybackState.UpdatedAt.UnixNano() / int64(time.Millisecond),
		"numOfParticipants": len(r.Participants),
	}

	msg, err := json.Marshal(roomState)
	if err != nil {
		return err
	}

	h.Broadcast(msg)
	return nil
}
