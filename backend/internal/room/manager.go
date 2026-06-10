package room

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/abram056/syncstream/backend/internal/models"
)

type Manager struct {
	repo Repository
}

func NewManager(repo Repository) *Manager {
	return &Manager{repo: repo}
}

var ErrRoomNotFound = errors.New("room not found")
var ErrRoomAlreadyExists = errors.New("room already exists")

func (m *Manager) CreateRoom(media models.Media) (*models.Room, error) {
	room := &models.Room{
		ID:           generateRoomID(),
		Status:       models.Waiting,
		Media:        media,
		Participants: make(map[string]*models.Participant),
		PlaybackState: models.PlaybackState{
			IsPlaying: false,
			Position:  0,
			UpdatedAt: time.Now(),
		},
		CreatedAt:    time.Now(),
		LastActiveAt: time.Now(),
	}

	if err := m.repo.CreateRoom(room); err != nil {
		return nil, ErrRoomAlreadyExists
	}

	return room, nil
}

func (m *Manager) GetRoomByID(id string) (*models.Room, error) {
	room, err := m.repo.GetRoomByID(id)
	if err != nil {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func generateRoomID() string {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().Format("rm20060102150405")
	}
	return "rm-" + hex.EncodeToString(buf)
}
