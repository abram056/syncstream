package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/abram056/syncstream/backend/internal/api"
	"github.com/abram056/syncstream/backend/internal/logger"
	"github.com/abram056/syncstream/backend/internal/room"
	memory "github.com/abram056/syncstream/backend/internal/storage/memory"
)

func main() {
	logger.Init()
	cfg := room.DefaultConfig()

	// create repository and manager shared across handlers and background tasks
	repo := memory.NewRoomStore()
	manager := room.NewManager(repo)

	srv := api.NewServer(manager)
	hubReg := srv.HubRegistry

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      srv.Handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// background participant expiration: periodically remove stale disconnected participants
	go func() {
		ticker := time.NewTicker(cfg.CleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			expired := manager.ExpireParticipants(cfg.ParticipantGracePeriod)
			for roomID, participantIDs := range expired {
				for _, pid := range participantIDs {
					slog.Info("expired participant",
						"participant_id", pid,
						"room_id", roomID,
					)
				}
				// Broadcast user_left for expired participants via the room's hub
				if hub := hubReg.Get(roomID); hub != nil {
					for _, pid := range participantIDs {
						evt := map[string]interface{}{
							"type":   "user_left",
							"userId": pid,
						}
						if msg, err := json.Marshal(evt); err == nil {
							hub.Broadcast(msg)
						} else {
							slog.Warn("failed to marshal user_left event for expired participant",
								"room_id", roomID,
								"participant_id", pid,
								"error", err,
							)
						}
					}
					hub.BroadcastRoomState("room_state")
				}
			}
		}
	}()

	// background idle room cleanup: remove inactive rooms after the configured timeout
	go func() {
		ticker := time.NewTicker(cfg.CleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			removed, err := manager.CleanupIdleRooms(cfg.RoomIdleTimeout)
			if err != nil {
				slog.Error("room cleanup error", "error", err)
				continue
			}
			for _, roomID := range removed {
				hubReg.Remove(roomID)
				slog.Info("cleaned up idle room", "room_id", roomID)
			}
		}
	}()

	slog.Info("starting server", "addr", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
	}
}
