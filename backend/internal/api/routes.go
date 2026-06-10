package api

import (
	"net/http"

	"github.com/abram056/syncstream/backend/internal/api/handlers"
)

func NewRouter() http.Handler {
	h := handlers.NewHandler()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/v1/rooms", h.CreateRoom)
	mux.HandleFunc("/api/v1/rooms/", h.GetRoom)

	return mux
}
