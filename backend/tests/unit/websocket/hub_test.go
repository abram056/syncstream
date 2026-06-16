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

// TestHubClientRegistration verifies clients can register with a hub
func TestHubClientRegistration(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())

	// create a test room
	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	// create mock clients
	client1 := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-test-001",
		DisplayName:   "User 1",
		Joined:        false,
	}

	client2 := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-test-002",
		DisplayName:   "User 2",
		Joined:        false,
	}

	// register clients
	hub.RegisterClient(client1)
	time.Sleep(10 * time.Millisecond) // allow time for registration

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}

	hub.RegisterClient(client2)
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 2 {
		t.Errorf("expected 2 clients, got %d", hub.ClientCount())
	}
}

// TestHubHandleJoinRoom verifies join_room event handling
func TestHubHandleJoinRoom(t *testing.T) {
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
		ParticipantID: "usr-test-join",
		DisplayName:   "",
		Joined:        false,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	// send join_room event
	evt := map[string]interface{}{
		"type":         "join_room",
		"display_name": "Test User",
	}

	err = hub.HandleJoinRoom(client, evt)
	if err != nil {
		t.Fatalf("failed to handle join_room: %v", err)
	}

	if client.DisplayName != "Test User" {
		t.Errorf("expected display name 'Test User', got %s", client.DisplayName)
	}

	// verify room_joined message was sent
	select {
	case msg := <-client.Send:
		var result map[string]interface{}
		if err := json.Unmarshal(msg, &result); err != nil {
			t.Fatalf("failed to unmarshal message: %v", err)
		}
		if result["type"] != "room_joined" {
			t.Errorf("expected room_joined message, got %v", result["type"])
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for room_joined message")
	}
}

// TestHubBroadcast verifies broadcast to all clients
func TestHubBroadcast(t *testing.T) {
	manager := room.NewManager(memory.NewRoomStore())

	r, err := manager.CreateRoom(models.Media{URL: "http://example.com/video.mp4", Title: "Test"})
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	hub := websocket.NewHub(r.ID, manager)
	go hub.Run()

	client1 := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-bc-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	client2 := &websocket.Client{
		Hub:           hub,
		Send:          make(chan []byte, 256),
		ParticipantID: "usr-bc-002",
		DisplayName:   "User 2",
		Joined:        true,
	}

	hub.RegisterClient(client1)
	hub.RegisterClient(client2)
	time.Sleep(10 * time.Millisecond)

	// broadcast a message
	msg := []byte(`{"type":"test_event","data":"hello"}`)
	hub.Broadcast(msg)

	// verify both clients receive it
	for _, client := range []*websocket.Client{client1, client2} {
		select {
		case received := <-client.Send:
			if string(received) != string(msg) {
				t.Errorf("unexpected message: %s", string(received))
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timeout waiting for broadcast message")
		}
	}
}

// TestHubClientUnregister verifies client unregistration
func TestHubClientUnregister(t *testing.T) {
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
		ParticipantID: "usr-unreg-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client after registration, got %d", hub.ClientCount())
	}

	// unregister client
	hub.UnregisterClient(client)
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients after unregistration, got %d", hub.ClientCount())
	}
}

// TestHubStop verifies that stopping a hub terminates its Run goroutine
// and closes all client connections.
func TestHubStop(t *testing.T) {
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
		ParticipantID: "usr-stop-001",
		DisplayName:   "User 1",
		Joined:        true,
	}

	hub.RegisterClient(client)
	time.Sleep(10 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client before stop, got %d", hub.ClientCount())
	}

	// Stop the hub
	hub.Stop()
	time.Sleep(10 * time.Millisecond)

	// After stop, the hub should have no clients
	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients after stop, got %d", hub.ClientCount())
	}

	// Sending to a stopped hub should not panic (broadcast channel should still be open
	// or safely handled)
	hub.Broadcast([]byte(`{"type":"test"}`))
}
