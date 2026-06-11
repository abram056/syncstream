package room

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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

// AddParticipant adds a participant to the room and updates room status.
func (m *Manager) AddParticipant(roomID, participantID, displayName string) (*models.Participant, error) {
	r, err := m.repo.GetRoomByID(roomID)
	if err != nil {
		return nil, ErrRoomNotFound
	}

	// if participant already exists, mark connected
	if p, ok := r.Participants[participantID]; ok {
		p.Connected = true
		p.JoinedAt = time.Now()
		r.LastActiveAt = time.Now()
		r.Status = models.Active
		return p, nil
	}

	p := &models.Participant{
		ID:          participantID,
		DisplayName: displayName,
		Connected:   true,
		JoinedAt:    time.Now(),
	}
	r.Participants[participantID] = p
	r.LastActiveAt = time.Now()
	r.Status = models.Active

	return p, nil
}

// RemoveParticipant removes a participant and updates room status.
func (m *Manager) RemoveParticipant(roomID, participantID string) error {
	r, err := m.repo.GetRoomByID(roomID)
	if err != nil {
		return ErrRoomNotFound
	}

	if _, ok := r.Participants[participantID]; !ok {
		return fmt.Errorf("participant %s not found", participantID)
	}

	delete(r.Participants, participantID)
	r.LastActiveAt = time.Now()

	if len(r.Participants) == 0 {
		r.Status = models.Idle
	}

	return nil
}

// CleanupIdleRooms removes rooms that are idle for longer than threshold.
func (m *Manager) CleanupIdleRooms(threshold time.Duration) ([]string, error) {
	rooms, err := m.repo.ListRooms()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var removed []string
	for _, r := range rooms {
		if r.Status == models.Idle && now.Sub(r.LastActiveAt) > threshold {
			if err := m.repo.DeleteRoom(r.ID); err == nil {
				removed = append(removed, r.ID)
			}
		}
	}
	return removed, nil
}
