package room_test

import (
	"testing"

	"github.com/abram056/syncstream/backend/internal/models"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

func TestCreateRoomSuccess(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	media := models.Media{URL: "https://example.com/video.mp4", Title: "Test Video"}
	roomObj, err := manager.CreateRoom(media)
	if err != nil {
		t.Fatalf("expected create room to succeed, got %v", err)
	}

	if roomObj.ID == "" {
		t.Fatal("expected generated room ID")
	}

	if roomObj.Media.URL != media.URL {
		t.Fatalf("expected media URL %q but got %q", media.URL, roomObj.Media.URL)
	}
}

func TestGetRoomByIDNotFound(t *testing.T) {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	_, err := manager.GetRoomByID("missing-room")
	if err == nil {
		t.Fatal("expected error when room is missing")
	}
}
