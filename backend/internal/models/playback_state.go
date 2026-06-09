package models

import "time"

type PlaybackState struct {
	IsPlaying bool
	Position  float64

	UpdatedBy string
	UpdatedAt time.Time
}
