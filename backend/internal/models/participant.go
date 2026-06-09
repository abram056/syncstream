package models

import "time"

type Participant struct {
	ID         string
	DsplayName string
	Connected  bool
	JoinedAt   time.Time
	CanControl bool
}
