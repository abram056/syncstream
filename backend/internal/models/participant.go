package models

import "time"

type Participant struct {
	ID          string
	DisplayName string
	Connected   bool
	JoinedAt    time.Time
	LastSeen    time.Time
	CanControl  bool
}
