package models

import (
	"sync"
	"time"
)

type RoomStatus string

const (
	Waiting RoomStatus = "waiting"
	Active  RoomStatus = "active"
	Idle    RoomStatus = "idle"
)

type Room struct {
	mu            sync.RWMutex
	ID            string
	Status        RoomStatus
	Media         Media
	Participants  map[string]*Participant
	PlaybackState PlaybackState
	CreatedAt     time.Time
	LastActiveAt  time.Time
}

func (r *Room) Lock()   { r.mu.Lock() }
func (r *Room) Unlock() { r.mu.Unlock() }

func (r *Room) RLock()   { r.mu.RLock() }
func (r *Room) RUnlock() { r.mu.RUnlock() }
