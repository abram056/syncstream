package storage_test

import (
	"testing"

	"github.com/abram056/syncstream/backend/internal/models"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

func TestRoomStoreCreateAndGet(t *testing.T) {
	store := memory.NewRoomStore()

	room := &models.Room{
		ID:    "rm-test",
		Media: models.Media{URL: "https://example.com/video.mp4", Title: "Test Video"},
	}

	if err := store.CreateRoom(room); err != nil {
		t.Fatalf("expected create room to succeed, got %v", err)
	}

	loaded, err := store.GetRoomByID("rm-test")
	if err != nil {
		t.Fatalf("expected get room to succeed, got %v", err)
	}

	if loaded.ID != room.ID {
		t.Fatalf("expected loaded room ID %q but got %q", room.ID, loaded.ID)
	}
}

func TestRoomStoreDuplicateRoom(t *testing.T) {
	store := memory.NewRoomStore()

	room := &models.Room{ID: "rm-test"}
	if err := store.CreateRoom(room); err != nil {
		t.Fatalf("expected create room to succeed, got %v", err)
	}

	if err := store.CreateRoom(room); err == nil {
		t.Fatal("expected duplicate create to fail")
	}
}

func TestRoomStoreDeleteRoom(t *testing.T) {
	store := memory.NewRoomStore()

	room := &models.Room{ID: "rm-test"}
	if err := store.CreateRoom(room); err != nil {
		t.Fatalf("expected create room to succeed, got %v", err)
	}

	if err := store.DeleteRoom("rm-test"); err != nil {
		t.Fatalf("expected delete room to succeed, got %v", err)
	}

	if _, err := store.GetRoomByID("rm-test"); err == nil {
		t.Fatal("expected get room after delete to fail")
	}
}
