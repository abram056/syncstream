package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abram056/syncstream/backend/internal/api"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

func main() {
	// create repository and manager shared across handlers and background tasks
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	router := api.NewRouter(manager)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// background cleanup ticker: remove idle rooms older than 1 hour every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			removed, err := manager.CleanupIdleRooms(1 * time.Hour)
			if err != nil {
				log.Printf("room cleanup error: %v", err)
				continue
			}
			if len(removed) > 0 {
				log.Printf("cleaned up rooms: %v", removed)
			}
		}
	}()

	log.Printf("starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
