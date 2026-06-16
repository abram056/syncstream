package websocket_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/abram056/syncstream/backend/internal/models"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
	"github.com/abram056/syncstream/backend/internal/websocket"
)

// TestHandlePlay verifies play event handling and sync_state broadcast
func TestHandlePlay(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	client := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-play-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	// Add participant to room
	_, err = manager.AddParticipant(r.ID, client.ParticipantID, client.DisplayName)
	if err != nil {
		t.Fatalf("failed to add participant: %v", err)
	}

	// Handle play event
	if err := hub.HandlePlay(client, 0, false); err != nil {
		t.Fatalf("failed to handle play: %v", err)
	}

	// Verify room playback state changed
	updated, err := manager.GetRoomByID(r.ID)
	if err != nil {
		t.Fatalf("failed to get room: %v", err)
	}

	if !updated.PlaybackState.IsPlaying {
		t.Error("expected IsPlaying to be true")
	}

	if updated.PlaybackState.UpdatedBy != client.ParticipantID {
		t.Errorf("expected UpdatedBy to be %s, got %s", client.ParticipantID, updated.PlaybackState.UpdatedBy)
	}

	// Verify sync_state was broadcast
	select {
	case msg := <-client.Send:
		var result map[string]interface{}
		if err := json.Unmarshal(msg, &result); err != nil {
			t.Fatalf("failed to unmarshal message: %v", err)
		}
		if result["type"] != "sync_state" {
			t.Errorf("expected sync_state message, got %v", result["type"])
		}
		if result["isPlaying"] != true {
			t.Error("expected isPlaying to be true in broadcast")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for sync_state broadcast")
	}
}

// TestHandlePause verifies pause event handling
func TestHandlePause(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	client := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-pause-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	_, err = manager.AddParticipant(r.ID, client.ParticipantID, client.DisplayName)
	if err != nil {
		t.Fatalf("failed to add participant: %v", err)
	}

	// First play
	if _, err := manager.Play(r.ID, client.ParticipantID); err != nil {
		t.Fatalf("failed to play: %v", err)
	}

	// Then pause
	if err := hub.HandlePause(client, 0, false); err != nil {
		t.Fatalf("failed to handle pause: %v", err)
	}

	// Verify room playback state changed
	updated, err := manager.GetRoomByID(r.ID)
	if err != nil {
		t.Fatalf("failed to get room: %v", err)
	}

	if updated.PlaybackState.IsPlaying {
		t.Error("expected IsPlaying to be false after pause")
	}
}

// TestHandleSeek verifies seek event handling with position
func TestHandleSeek(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	client := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-seek-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	_, err = manager.AddParticipant(r.ID, client.ParticipantID, client.DisplayName)
	if err != nil {
		t.Fatalf("failed to add participant: %v", err)
	}

	// Handle seek event with position
	if err := hub.HandleSeek(client, 42.5); err != nil {
		t.Fatalf("failed to handle seek: %v", err)
	}

	// Verify position was updated
	updated, err := manager.GetRoomByID(r.ID)
	if err != nil {
		t.Fatalf("failed to get room: %v", err)
	}

	if updated.PlaybackState.Position != 42.5 {
		t.Errorf("expected position 42.5, got %v", updated.PlaybackState.Position)
	}
}

// TestHandlePlayWithPosition verifies play with position updates both IsPlaying and Position
func TestHandlePlayWithPosition(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	client := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-play-pos-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	_, err = manager.AddParticipant(r.ID, client.ParticipantID, client.DisplayName)
	if err != nil {
		t.Fatalf("failed to add participant: %v", err)
	}

	// Handle play event with position
	if err := hub.HandlePlay(client, 15.0, true); err != nil {
		t.Fatalf("failed to handle play with position: %v", err)
	}

	// Verify both position and playback state
	updated, err := manager.GetRoomByID(r.ID)
	if err != nil {
		t.Fatalf("failed to get room: %v", err)
	}

	if !updated.PlaybackState.IsPlaying {
		t.Error("expected IsPlaying to be true")
	}

	if updated.PlaybackState.Position != 15.0 {
		t.Errorf("expected position 15.0, got %v", updated.PlaybackState.Position)
	}
}

// TestInvalidPositionRejected verifies negative positions are rejected
func TestInvalidPositionRejected(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	_, err = manager.AddParticipant(r.ID, "usr-invalid-001", "User 1")
	if err != nil {
		t.Fatalf("failed to add participant: %v", err)
	}

	// Try to seek to negative position
	_, err = manager.Seek(r.ID, "usr-invalid-001", -5.0)
	if err == nil {
		t.Error("expected error for negative position")
	}
}

// TestSyncStateBroadcastIncludesMetadata verifies sync_state contains all required fields
func TestSyncStateBroadcastIncludesMetadata(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	client := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-meta-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	_, err = manager.AddParticipant(r.ID, client.ParticipantID, client.DisplayName)
	if err != nil {
		t.Fatalf("failed to add participant: %v", err)
	}

	if err := hub.HandlePlay(client, 10.0, true); err != nil {
		t.Fatalf("failed to handle play: %v", err)
	}

	// Collect broadcast
	select {
	case msg := <-client.Send:
		var result map[string]interface{}
		if err := json.Unmarshal(msg, &result); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		// Verify all required fields
		requiredFields := []string{"type", "roomId", "status", "isPlaying", "position", "updatedBy", "updatedAt", "numOfParticipants"}
		for _, field := range requiredFields {
			if _, ok := result[field]; !ok {
				t.Errorf("missing field: %s", field)
			}
		}

		if result["isPlaying"] != true {
			t.Error("expected isPlaying to be true")
		}
		pos, ok := result["position"].(float64)
		if !ok {
			t.Error("expected position to be a float64")
		} else if pos < 10.0 {
			t.Errorf("expected position >= 10.0, got %v", pos)
		}
		if result["updatedBy"] != client.ParticipantID {
			t.Errorf("expected updatedBy to be %s", client.ParticipantID)
		}

	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for sync_state")
	}
}

// TestEffectivePosition verifies the authoritative position calculation:
// - Paused playback returns the stored position unchanged
// - Active playback returns a position that advances with elapsed time
func TestEffectivePosition(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	// When paused, effective position equals stored position
	r.PlaybackState.IsPlaying = false
	r.PlaybackState.Position = 42.0
	r.PlaybackState.UpdatedAt = time.Now().Add(-5 * time.Second)

	effectivePaused := r.PlaybackState.EffectivePosition()
	if effectivePaused != 42.0 {
		t.Errorf("expected paused effective position 42.0, got %v", effectivePaused)
	}

	// When playing, effective position advances with time
	r.PlaybackState.IsPlaying = true
	r.PlaybackState.Position = 100.0
	r.PlaybackState.UpdatedAt = time.Now().Add(-2 * time.Second)

	effectivePlaying := r.PlaybackState.EffectivePosition()
	if effectivePlaying < 100.0 {
		t.Errorf("expected playing effective position >= 100.0, got %v", effectivePlaying)
	}
	if effectivePlaying > 200.0 {
		t.Errorf("expected playing effective position <= 200.0 (sanity), got %v", effectivePlaying)
	}
}
