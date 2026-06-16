package room

import "time"

type Config struct {
	ParticipantGracePeriod time.Duration
	RoomIdleTimeout        time.Duration
	CleanupInterval        time.Duration
}

func DefaultConfig() Config {
	return Config{
		ParticipantGracePeriod: 5 * time.Minute,
		RoomIdleTimeout:        30 * time.Minute,
		CleanupInterval:        10 * time.Minute,
	}
}
