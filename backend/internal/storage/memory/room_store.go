package storage

import (
	"fmt"
	"sync"

	"github.com/abram056/syncstream/backend/internal/models"
)

type RoomStore struct {
	rooms map[string]*models.Room
	mu    sync.RWMutex
}

func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[string]*models.Room),
	}
}

func (s *RoomStore) CreateRoom(room *models.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.rooms[room.ID]; exists {
		return fmt.Errorf("roomID %s already exists", room.ID)
	}

	s.rooms[room.ID] = room
	return nil
}

func (s *RoomStore) GetRoomByID(id string) (*models.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[id]
	if !exists {
		return nil, fmt.Errorf("roomID %s not found", id)
	}

	return room, nil
}

func (s *RoomStore) DeleteRoom(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.rooms[id]; !exists {
		return fmt.Errorf("roomID %s not found", id)
	}

	delete(s.rooms, id)
	return nil
}
