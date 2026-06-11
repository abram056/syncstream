package api

import (
	"net/http"

	"github.com/abram056/syncstream/backend/internal/api/handlers"
	"github.com/abram056/syncstream/backend/internal/room"
)

func NewRouter(manager *room.Manager) http.Handler {
	h := handlers.NewHandler(manager)
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/v1/rooms", h.CreateRoom)
	mux.HandleFunc("/api/v1/rooms/", h.GetRoom)

	return mux
}
