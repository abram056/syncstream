package websocket

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/abram056/syncstream/backend/internal/room"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for MVP
		},
	}
	hubs  = make(map[string]*Hub)
	hubMu sync.RWMutex
)

// Handler handles WebSocket upgrade requests and room connections.
type Handler struct {
	manager *room.Manager
}

// NewHandler creates a new WebSocket handler.
func NewHandler(manager *room.Manager) *Handler {
	return &Handler{manager: manager}
}

// ServeWS handles WebSocket connections.
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Extract room_id from URL path: /api/v1/rooms/{room_id}/ws
	const prefix = "/api/v1/rooms/"
	const suffix = "/ws"

	if !strings.HasPrefix(r.URL.Path, prefix) || !strings.HasSuffix(r.URL.Path, suffix) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Extract room_id
	roomID := strings.TrimPrefix(r.URL.Path, prefix)
	roomID = strings.TrimSuffix(roomID, suffix)

	if roomID == "" {
		http.Error(w, "missing room_id", http.StatusBadRequest)
		return
	}

	// Verify room exists
	if _, err := h.manager.GetRoomByID(roomID); err != nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	// Get or create hub for this room
	hub := getOrCreateHub(roomID, h.manager)

	// Create client and register
	client := NewClient(hub, conn)
	hub.register <- client

	// Start read/write pumps
	go client.WritePump()
	go client.ReadPump()
}

// getOrCreateHub gets an existing hub or creates a new one for a room.
func getOrCreateHub(roomID string, manager *room.Manager) *Hub {
	hubMu.Lock()
	defer hubMu.Unlock()

	if hub, ok := hubs[roomID]; ok {
		return hub
	}

	hub := NewHub(roomID, manager)
	hubs[roomID] = hub
	go hub.Run()
	return hub
}
