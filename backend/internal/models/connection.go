package models

import "time"

type Connection struct {
	ParticipantID string

	ConnectedAt time.Time
	LastPingAt  time.Time
}
