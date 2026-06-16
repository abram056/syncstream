package websocket

import (
	"sync"

	"github.com/abram056/syncstream/backend/internal/room"
)

// HubRegistry manages the lifecycle of WebSocket hubs indexed by room ID.
// It provides thread-safe access for creating, retrieving, and removing hubs.
type HubRegistry struct {
	mu   sync.RWMutex
	hubs map[string]*Hub
}

func NewHubRegistry() *HubRegistry {
	return &HubRegistry{
		hubs: make(map[string]*Hub),
	}
}

func (reg *HubRegistry) GetOrCreate(roomID string, manager *room.Manager) *Hub {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if hub, ok := reg.hubs[roomID]; ok {
		return hub
	}

	hub := NewHub(roomID, manager)
	reg.hubs[roomID] = hub
	go hub.Run()
	return hub
}

func (reg *HubRegistry) Get(roomID string) *Hub {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.hubs[roomID]
}

func (reg *HubRegistry) Remove(roomID string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if hub, ok := reg.hubs[roomID]; ok {
		hub.Stop()
		delete(reg.hubs, roomID)
	}
}

func (reg *HubRegistry) List() []*Hub {
	reg.mu.RLock()
	defer reg.mu.RUnlock()

	hubs := make([]*Hub, 0, len(reg.hubs))
	for _, hub := range reg.hubs {
		hubs = append(hubs, hub)
	}
	return hubs
}
