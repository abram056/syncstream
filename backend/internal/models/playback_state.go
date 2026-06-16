package models

import "time"

type PlaybackState struct {
	IsPlaying bool
	Position  float64

	UpdatedBy string
	UpdatedAt time.Time
}

// EffectivePosition calculates the authoritative playback position.
// If playback is active, the position advances by the elapsed time since
// the last update. This ensures reconnecting users and new joiners receive
// an accurate position rather than a stale snapshot.
func (s *PlaybackState) EffectivePosition() float64 {
	if !s.IsPlaying {
		return s.Position
	}
	elapsed := time.Since(s.UpdatedAt).Seconds()
	if elapsed < 0 {
		return s.Position
	}
	return s.Position + elapsed
}
