package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abram056/syncstream/backend/internal/models"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

type Handler struct {
	roomManager *room.Manager
}

type createRoomRequest struct {
	MediaURL string `json:"media_url"`
	Title    string `json:"title,omitempty"`
}

type createRoomResponse struct {
	RoomID string `json:"room_id"`
}

type getRoomResponse struct {
	RoomID       string  `json:"room_id"`
	Status       string  `json:"status"`
	MediaURL     string  `json:"media_url"`
	Title        string  `json:"title,omitempty"`
	IsPlaying    bool    `json:"is_playing"`
	Position     float64 `json:"position"`
	Participants int     `json:"participants"`
}

func NewHandler() *Handler {
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)
	return &Handler{roomManager: manager}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request payload"})
		return
	}

	if req.MediaURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "media_url is required"})
		return
	}

	room, err := h.roomManager.CreateRoom(models.Media{URL: req.MediaURL, Title: req.Title})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unable to create room"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createRoomResponse{RoomID: room.ID})
}

func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	const prefix = "/api/v1/rooms/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "room not found"})
		return
	}

	roomID := strings.TrimPrefix(r.URL.Path, prefix)
	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "room_id is required"})
		return
	}

	roomData, err := h.roomManager.GetRoomByID(roomID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "room not found"})
		return
	}

	resp := getRoomResponse{
		RoomID:       roomData.ID,
		Status:       string(roomData.Status),
		MediaURL:     roomData.Media.URL,
		Title:        roomData.Media.Title,
		IsPlaying:    roomData.PlaybackState.IsPlaying,
		Position:     roomData.PlaybackState.Position,
		Participants: len(roomData.Participants),
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
