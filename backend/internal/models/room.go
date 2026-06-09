package models

import "time"

type RoomStatus string

const (
	Waiting RoomStatus = "waiting"
	Active  RoomStatus = "active"
	Idle    RoomStatus = "idle"
)

type Room struct {
	ID            string
	Status        RoomStatus
	Media         Media
	Participants  map[string]*Participant
	PlaybackState PlaybackState
	CreatedAt     time.Time
	LastActiveAt  time.Time
}
