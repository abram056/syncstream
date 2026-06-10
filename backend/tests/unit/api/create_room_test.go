package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abram056/syncstream/backend/internal/api"
)

type createRoomResponse struct {
	RoomID string `json:"room_id"`
}

func TestCreateRoomEndpoint(t *testing.T) {
	router := api.NewRouter()

	payload := map[string]string{"media_url": "https://example.com/video.mp4", "title": "Test Video"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewReader(body))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d but got %d", http.StatusCreated, res.Code)
	}

	var resp createRoomResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if resp.RoomID == "" {
		t.Fatal("expected non-empty room_id")
	}
}

func TestGetRoomEndpoint(t *testing.T) {
	router := api.NewRouter()

	// First create a room.
	payload := map[string]string{"media_url": "https://example.com/video.mp4"}
	body, _ := json.Marshal(payload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewReader(body))
	createRes := httptest.NewRecorder()
	router.ServeHTTP(createRes, createReq)

	if createRes.Code != http.StatusCreated {
		t.Fatalf("expected create status %d but got %d", http.StatusCreated, createRes.Code)
	}

	var createResp createRoomResponse
	if err := json.NewDecoder(createRes.Body).Decode(&createResp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+createResp.RoomID, nil)
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)

	if getRes.Code != http.StatusOK {
		t.Fatalf("expected get status %d but got %d", http.StatusOK, getRes.Code)
	}
}
