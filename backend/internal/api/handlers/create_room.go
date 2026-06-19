package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/abram056/syncstream/backend/internal/models"
	"github.com/abram056/syncstream/backend/internal/room"
	ws "github.com/abram056/syncstream/backend/internal/websocket"
)

type Handler struct {
	roomManager *room.Manager
	wsHandler   *ws.Handler
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

func NewHandler(manager *room.Manager) *Handler {
	return &Handler{
		roomManager: manager,
		wsHandler:   ws.NewHandler(manager),
	}
}

func (h *Handler) HubRegistry() *ws.HubRegistry {
	return h.wsHandler.HubRegistry()
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		slog.Warn("Health: failed to encode response", "error", err)
	}
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"}); err != nil {
			slog.Warn("CreateRoom: failed to encode response", "error", err)
		}
		slog.Warn("CreateRoom: method not allowed", "method", r.Method)
		return
	}

	var req createRoomRequest
	slog.Debug("CreateRoom: request", "uri", r.RequestURI)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "invalid request payload"}); err != nil {
			slog.Warn("CreateRoom: failed to encode response", "error", err)
		}
		slog.Warn("CreateRoom: failed to decode request body", "error", err)
		return
	}

	if req.MediaURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "media_url is required"}); err != nil {
			slog.Warn("CreateRoom: failed to encode response", "error", err)
		}
		slog.Warn("CreateRoom: media_url is required")
		return
	}

	createdRoom, err := h.roomManager.CreateRoom(models.Media{URL: req.MediaURL, Title: req.Title})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "unable to create room"}); err != nil {
			slog.Warn("CreateRoom: failed to encode response", "error", err)
		}
		slog.Error("CreateRoom: failed to create room", "error", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createRoomResponse{RoomID: createdRoom.ID}); err != nil {
		slog.Warn("CreateRoom: failed to encode response", "error", err)
	}
	slog.Info("CreateRoom: room created", "room_id", createdRoom.ID)
}

func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"}); err != nil {
			slog.Warn("GetRoom: failed to encode response", "error", err)
		}
		slog.Warn("GetRoom: method not allowed", "method", r.Method)
		return
	}

	const prefix = "/api/v1/rooms/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "room not found"}); err != nil {
			slog.Warn("GetRoom: failed to encode response", "error", err)
		}
		slog.Warn("GetRoom: room not found")
		return
	}

	roomID := strings.TrimPrefix(r.URL.Path, prefix)
	slog.Debug("GetRoom: fetching room", "room_id", roomID)
	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "room_id is required"}); err != nil {
			slog.Warn("GetRoom: failed to encode response", "error", err)
		}
		slog.Warn("GetRoom: room_id is required")
		return
	}

	roomData, err := h.roomManager.GetRoomByID(roomID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "room not found"}); err != nil {
			slog.Warn("GetRoom: failed to encode response", "error", err)
		}
		slog.Warn("GetRoom: failed to get room", "room_id", roomID, "error", err)
		return
	}

	resp := getRoomResponse{
		RoomID:       roomData.ID,
		Status:       string(roomData.Status),
		MediaURL:     roomData.Media.URL,
		Title:        roomData.Media.Title,
		IsPlaying:    roomData.PlaybackState.IsPlaying,
		Position:     roomData.PlaybackState.EffectivePosition(),
		Participants: len(roomData.Participants),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Warn("GetRoom: failed to encode response", "error", err)
	}
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	h.wsHandler.ServeWS(w, r)
}
