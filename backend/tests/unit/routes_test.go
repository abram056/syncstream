package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abram056/syncstream/backend/internal/api"
)

type healthResponse struct {
	Status string `json:"status"`
}

func TestHealthRoute(t *testing.T) {
	router := api.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, res.Code)
	}

	var payload healthResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Status != "ok" {
		t.Fatalf("expected status=ok but got %q", payload.Status)
	}
}

func TestCreateRoomMethodNotAllowed(t *testing.T) {
	router := api.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d but got %d", http.StatusMethodNotAllowed, res.Code)
	}
}

func TestCreateRoomNotImplemented(t *testing.T) {
	router := api.NewRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotImplemented {
		t.Fatalf("expected status %d but got %d", http.StatusNotImplemented, res.Code)
	}
}
