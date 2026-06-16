package room

import (
	"testing"
	"time"

	"github.com/abram056/syncstream/backend/internal/models"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

func TestCleanupIdleRooms(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	// Create idle room (last active 2 hours ago)
	r1, err := manager.CreateRoom(models.Media{URL: "https://example.com/idle.mp4"})
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}
	r1.Status = models.Idle
	r1.LastActiveAt = time.Now().Add(-2 * time.Hour)

	// Create active room (last active 2 hours ago but status active)
	r2, err := manager.CreateRoom(models.Media{URL: "https://example.com/active.mp4"})
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}
	r2.Status = models.Active
	r2.LastActiveAt = time.Now().Add(-2 * time.Hour)

	removed, err := manager.CleanupIdleRooms(1 * time.Hour)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	if len(removed) != 1 {
		t.Fatalf("expected 1 room removed, got %d", len(removed))
	}

	if _, err := repo.GetRoomByID(r1.ID); err == nil {
		t.Fatalf("expected idle room to be deleted")
	}

	if _, err := repo.GetRoomByID(r2.ID); err != nil {
		t.Fatalf("expected active room to remain, got error %v", err)
	}
}

// TestExpireParticipants verifies that disconnected participants past the grace
// period are removed, while connected participants and recently disconnected ones
// are preserved.
func TestExpireParticipants(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	r, err := manager.CreateRoom(models.Media{URL: "https://example.com/video.mp4"})
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}

	p1, _ := manager.AddParticipant(r.ID, "user1", "Alice")
	p2, _ := manager.AddParticipant(r.ID, "user2", "Bob")
	_, _ = manager.AddParticipant(r.ID, "user3", "Charlie")

	// Disconnect user1 — expired (set LastSeen far in the past)
	p1.Connected = false
	p1.LastSeen = time.Now().Add(-10 * time.Minute)

	// Disconnect user2 — recently disconnected, should NOT be expired
	p2.Connected = false
	p2.LastSeen = time.Now().Add(-1 * time.Minute)

	// user3 stays connected

	expired := manager.ExpireParticipants(5 * time.Minute)

	if len(expired[r.ID]) != 1 {
		t.Fatalf("expected 1 expired participant, got %d", len(expired[r.ID]))
	}

	if expired[r.ID][0] != "user1" {
		t.Fatalf("expected user1 to be expired, got %s", expired[r.ID][0])
	}

	// Verify user1 removed, others remain
	if _, ok := r.Participants["user1"]; ok {
		t.Fatal("expected user1 to be removed")
	}
	if _, ok := r.Participants["user2"]; !ok {
		t.Fatal("expected user2 to still be present")
	}
	if _, ok := r.Participants["user3"]; !ok {
		t.Fatal("expected user3 to still be present")
	}
}
