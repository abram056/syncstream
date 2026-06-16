package room_test

import (
	"testing"

	"github.com/abram056/syncstream/backend/internal/models"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

func TestParticipantJoinLeave(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	media := models.Media{URL: "https://example.com/video.mp4", Title: "Test"}
	r, err := manager.CreateRoom(media)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	p1, err := manager.AddParticipant(r.ID, "user1", "Alice")
	if err != nil {
		t.Fatalf("add participant failed: %v", err)
	}

	if !p1.Connected || p1.DisplayName != "Alice" {
		t.Fatalf("participant properties incorrect")
	}

	if r.Status != models.Active {
		t.Fatalf("expected room status active after join, got %v", r.Status)
	}

	p2, err := manager.AddParticipant(r.ID, "user2", "Bob")
	if err != nil {
		t.Fatalf("add participant 2 failed: %v", err)
	}

	if p2.CanControl {
		t.Fatalf("second participant should not have control")
	}

	// remove participants
	if err := manager.RemoveParticipant(r.ID, "user1"); err != nil {
		t.Fatalf("remove participant failed: %v", err)
	}

	if _, ok := r.Participants["user1"]; ok {
		t.Fatalf("participant user1 should be removed")
	}

	if r.Status != models.Active { // still active because user2 is present
		t.Fatalf("expected room status active after removing one participant, got %v", r.Status)
	}

	if err := manager.RemoveParticipant(r.ID, "user2"); err != nil {
		t.Fatalf("remove participant 2 failed: %v", err)
	}

	if r.Status != models.Idle {
		t.Fatalf("expected room status idle after removing all participants, got %v", r.Status)
	}
}

// TestDisconnectParticipant verifies that disconnecting a participant marks them
// as disconnected but does not remove them from the room.
func TestDisconnectParticipant(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	r, err := manager.CreateRoom(models.Media{URL: "https://example.com/video.mp4"})
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	p, err := manager.AddParticipant(r.ID, "user1", "Alice")
	if err != nil {
		t.Fatalf("add participant failed: %v", err)
	}

	if !p.Connected {
		t.Fatal("expected participant to be connected initially")
	}

	if err := manager.DisconnectParticipant(r.ID, "user1"); err != nil {
		t.Fatalf("disconnect participant failed: %v", err)
	}

	if p.Connected {
		t.Fatal("expected participant to be marked as disconnected")
	}

	if _, ok := r.Participants["user1"]; !ok {
		t.Fatal("expected participant to remain in room after disconnect")
	}

	if len(r.Participants) != 1 {
		t.Fatalf("expected participant count to remain 1, got %d", len(r.Participants))
	}
}

// TestReconnectParticipant verifies that a disconnected participant can reconnect
// using a valid reconnect token, and that an invalid token is rejected.
func TestReconnectParticipant(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	r, err := manager.CreateRoom(models.Media{URL: "https://example.com/video.mp4"})
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	p, err := manager.AddParticipant(r.ID, "user1", "Alice")
	if err != nil {
		t.Fatalf("add participant failed: %v", err)
	}

	if p.ReconnectToken == "" {
		t.Fatal("expected participant to have a reconnect token")
	}

	token := p.ReconnectToken

	// Disconnect the participant
	if err := manager.DisconnectParticipant(r.ID, "user1"); err != nil {
		t.Fatalf("disconnect participant failed: %v", err)
	}

	// Attempt reconnect with invalid token
	_, err = manager.ReconnectParticipant(r.ID, "user1", "invalid-token")
	if err != room.ErrReconnectTokenInvalid {
		t.Fatalf("expected ErrReconnectTokenInvalid, got %v", err)
	}

	// Attempt reconnect with valid token
	reconnected, err := manager.ReconnectParticipant(r.ID, "user1", token)
	if err != nil {
		t.Fatalf("reconnect with valid token failed: %v", err)
	}

	if !reconnected.Connected {
		t.Fatal("expected participant to be reconnected")
	}

	if reconnected.DisplayName != "Alice" {
		t.Fatalf("expected display name Alice, got %s", reconnected.DisplayName)
	}
}
