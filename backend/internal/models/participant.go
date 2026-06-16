package models

import "time"

type Participant struct {
	ID             string
	DisplayName    string
	Connected      bool
	JoinedAt       time.Time
	LastSeen       time.Time
	ReconnectToken string
	CanControl     bool
}

func (p *Participant) IsExpired(gracePeriod time.Duration) bool {
	if p.Connected {
		return false
	}
	return time.Since(p.LastSeen) > gracePeriod
}
