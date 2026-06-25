package websocket

import (
	"encoding/json"
	"log/slog"
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
	done       chan struct{}
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
		done:       make(chan struct{}),
	}
}

// Stop signals the hub's Run goroutine to shut down cleanly.
// It closes all client connections and releases hub resources.
func (h *Hub) Stop() {
	close(h.done)
}

// Run starts the hub's event loop.
// The loop exits when the done channel is closed (via Stop()).
// All client connections are cleaned up on exit.
func (h *Hub) Run() {
	for {
		select {
		case <-h.done:
			// Hub is shutting down. Close all client connections.
			h.mu.Lock()
			for client := range h.clients {
				close(client.Send)
				if client.Conn != nil {
					client.Conn.Close()
				}
			}
			h.clients = make(map[*Client]bool)
			h.mu.Unlock()
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			slog.Info("client registered",
				"room_id", h.roomID,
				"total_clients", len(h.clients),
			)

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
						slog.Warn("failed to disconnect participant",
							"room_id", h.roomID,
							"participant_id", client.ParticipantID,
							"error", err,
						)
					}

					// broadcast user_disconnected to remaining clients
					evt := map[string]interface{}{
						"type":        "user_disconnected",
						"userId":      client.ParticipantID,
						"displayName": client.DisplayName,
					}
					if msg, err := json.Marshal(evt); err == nil {
						h.Broadcast(msg)
					} else {
						slog.Warn("failed to marshal user_disconnected event",
							"room_id", h.roomID,
							"error", err,
						)
					}

					if err := h.broadcastRoomState("room_state"); err != nil {
						slog.Warn("failed to broadcast room_state after disconnect",
							"room_id", h.roomID,
							"error", err,
						)
					}
				}

				slog.Info("client unregistered",
					"room_id", h.roomID,
					"total_clients", len(h.clients),
				)
			} else {
				h.mu.Unlock()
			}

		case msg := <-h.broadcast:
			// log outgoing broadcast event type at debug level
			var evtType struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal(msg, &evtType); err == nil {
				slog.Debug("broadcasting event",
					"room_id", h.roomID,
					"event_type", evtType.Type,
				)
			}

			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- msg:
				default:
					slog.Warn("client send channel full, unregistering",
						"room_id", h.roomID,
						"participant_id", client.ParticipantID,
					)
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
// If the event includes a reconnect_token and participant_id, the server attempts
// to restore the previous session rather than creating a new participant.
// This allows clients to survive temporary disconnects without losing their identity.
func (h *Hub) HandleJoinRoom(client *Client, evt map[string]interface{}) error {
	displayName, ok := getDisplayName(evt)
	if !ok || displayName == "" {
		return ErrInvalidEvent
	}

	// Check for reconnect attempt
	reconnectToken, _ := evt["reconnect_token"].(string)
	if reconnectToken == "" {
		reconnectToken, _ = evt["reconnectToken"].(string)
	}

	// Use participant_id from event if provided (for cross-session reconnects)
	reconnectID := client.ParticipantID
	if pid, ok := evt["participant_id"].(string); ok && pid != "" {
		reconnectID = pid
	}

	var wasReconnect bool
	var reconnectTokenValue string

	if reconnectToken != "" && reconnectID != "" {
		p, err := h.manager.ReconnectParticipant(h.roomID, reconnectID, reconnectToken)
		if err == nil {
			client.DisplayName = p.DisplayName
			client.ParticipantID = p.ID
			reconnectTokenValue = p.ReconnectToken
			wasReconnect = true
		}
	}

	if !wasReconnect {
		p, err := h.manager.AddParticipant(h.roomID, client.ParticipantID, displayName)
		if err != nil {
			return err
		}
		client.DisplayName = p.DisplayName
		reconnectTokenValue = p.ReconnectToken
	}

	roomJoined := map[string]interface{}{
		"type":            "room_joined",
		"roomId":          h.roomID,
		"participantId":   client.ParticipantID,
		"reconnect_token": reconnectTokenValue,
	}
	if msg, err := json.Marshal(roomJoined); err == nil {
		client.Send <- msg
	} else {
		slog.Warn("failed to marshal room_joined event",
			"room_id", h.roomID,
			"error", err,
		)
	}

	r, err := h.manager.GetRoomByID(h.roomID)
	if err != nil {
		return err
	}

	effectivePos := r.PlaybackState.EffectivePosition()

	participantList := make([]map[string]interface{}, 0, len(r.Participants))
	for _, p := range r.Participants {
		participantList = append(participantList, map[string]interface{}{
			"userId":      p.ID,
			"displayName": p.DisplayName,
			"connected":   p.Connected,
		})
	}

	roomState := map[string]interface{}{
		"type":              "room_state",
		"roomId":            r.ID,
		"status":            string(r.Status),
		"mediaUrl":          r.Media.URL,
		"isPlaying":         r.PlaybackState.IsPlaying,
		"position":          effectivePos,
		"numOfParticipants": len(r.Participants),
		"participants":      participantList,
	}
	if msg, err := json.Marshal(roomState); err == nil {
		client.Send <- msg
	} else {
		slog.Warn("failed to marshal room_state event",
			"room_id", h.roomID,
			"error", err,
		)
	}

	// Broadcast appropriate event based on whether this is a new join or a reconnect
	if wasReconnect {
		userReconnected := map[string]interface{}{
			"type":        "user_reconnected",
			"userId":      client.ParticipantID,
			"displayName": displayName,
		}
		if msg, err := json.Marshal(userReconnected); err == nil {
			h.Broadcast(msg)
		} else {
			slog.Warn("failed to marshal user_reconnected event",
				"room_id", h.roomID,
				"error", err,
			)
		}
	} else {
		userJoined := map[string]interface{}{
			"type":        "user_joined",
			"userId":      client.ParticipantID,
			"displayName": displayName,
		}
		if msg, err := json.Marshal(userJoined); err == nil {
			h.Broadcast(msg)
		} else {
			slog.Warn("failed to marshal user_joined event",
				"room_id", h.roomID,
				"error", err,
			)
		}
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

// HandleLeaveRoom permanently removes a participant from the room and broadcasts user_left.
// This is distinct from a websocket disconnect, which only marks the participant as disconnected
// and broadcasts user_disconnected. A leave is an explicit intent to depart.
func (h *Hub) HandleLeaveRoom(client *Client) error {
	if client.ParticipantID == "" {
		return nil
	}

	if err := h.manager.RemoveParticipant(h.roomID, client.ParticipantID); err != nil {
		return err
	}

	evt := map[string]interface{}{
		"type":        "user_left",
		"userId":      client.ParticipantID,
		"displayName": client.DisplayName,
	}
	if msg, err := json.Marshal(evt); err == nil {
		h.Broadcast(msg)
	} else {
		slog.Warn("failed to marshal user_left event",
			"room_id", h.roomID,
			"error", err,
		)
	}

	return h.broadcastRoomState("room_state")
}

// BroadcastRoomState fetches the current room state and broadcasts it to all
// connected clients. This is used by cleanup services to notify participants
// of state changes (e.g., after participant expiration).
func (h *Hub) BroadcastRoomState(eventType string) error {
	return h.broadcastRoomState(eventType)
}

func (h *Hub) broadcastRoomState(eventType string) error {
	r, err := h.manager.GetRoomByID(h.roomID)
	if err != nil {
		return err
	}

	effectivePos := r.PlaybackState.EffectivePosition()

	roomState := map[string]interface{}{
		"type":              eventType,
		"roomId":            r.ID,
		"status":            string(r.Status),
		"mediaUrl":          r.Media.URL,
		"isPlaying":         r.PlaybackState.IsPlaying,
		"position":          effectivePos,
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
