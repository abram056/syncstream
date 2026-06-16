package api

import (
	"net/http"
	"strings"

	"github.com/abram056/syncstream/backend/internal/api/handlers"
	"github.com/abram056/syncstream/backend/internal/room"
	"github.com/abram056/syncstream/backend/internal/websocket"
)

type Server struct {
	Handler    http.Handler
	HubRegistry *websocket.HubRegistry
}

func NewServer(manager *room.Manager) *Server {
	h := handlers.NewHandler(manager)
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/v1/rooms", h.CreateRoom)
	mux.HandleFunc("/api/v1/rooms/", routeRoomRequest(h))

	return &Server{
		Handler:     mux,
		HubRegistry: h.HubRegistry(),
	}
}

// routeRoomRequest routes requests to either GetRoom or WebSocket based on the path.
func routeRoomRequest(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/ws") {
			h.WebSocket(w, r)
		} else {
			h.GetRoom(w, r)
		}
	}
}
