package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexmented/ai-for-developers-project-386/backend/internal/api"
	"github.com/alexmented/ai-for-developers-project-386/backend/internal/service"
)

func TestGetOwnerProfile(t *testing.T) {
	now := time.Date(2026, 3, 31, 6, 0, 0, 0, time.UTC)
	svc := service.NewCalendarService(func() time.Time { return now })
	server := NewServer(svc)
	handler := api.Handler(server)

	req := httptest.NewRequest(http.MethodGet, "/public/name-owner", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var profile api.PublicOwnerProfile
	if err := json.Unmarshal(recorder.Body.Bytes(), &profile); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if profile.Slug != "name-owner" {
		t.Fatalf("expected name-owner, got %s", profile.Slug)
	}
}

func TestCreateBookingConflictReturns409(t *testing.T) {
	now := time.Date(2026, 3, 31, 6, 0, 0, 0, time.UTC)
	svc := service.NewCalendarService(func() time.Time { return now })
	server := NewServer(svc)
	handler := api.Handler(server)

	firstPayload := map[string]any{
		"eventTypeId": "meeting-15",
		"slotStartAt": now.Format(time.RFC3339),
		"guestName":   "Ivan",
		"guestEmail":  "ivan@example.com",
	}

	firstBody, _ := json.Marshal(firstPayload)
	firstReq := httptest.NewRequest(http.MethodPost, "/public/name-owner/bookings", bytes.NewReader(firstBody))
	firstReq.Header.Set("Content-Type", "application/json")
	firstRecorder := httptest.NewRecorder()
	handler.ServeHTTP(firstRecorder, firstReq)

	if firstRecorder.Code != http.StatusOK {
		t.Fatalf("expected first booking 200, got %d", firstRecorder.Code)
	}

	secondPayload := map[string]any{
		"eventTypeId": "meeting-30",
		"slotStartAt": now.Format(time.RFC3339),
		"guestName":   "Petr",
		"guestEmail":  "petr@example.com",
	}

	secondBody, _ := json.Marshal(secondPayload)
	secondReq := httptest.NewRequest(http.MethodPost, "/public/name-owner/bookings", bytes.NewReader(secondBody))
	secondReq.Header.Set("Content-Type", "application/json")
	secondRecorder := httptest.NewRecorder()
	handler.ServeHTTP(secondRecorder, secondReq)

	if secondRecorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", secondRecorder.Code)
	}

	var apiErr api.ErrorResponse
	if err := json.Unmarshal(secondRecorder.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("failed to parse error body: %v", err)
	}

	if apiErr.Code != "SLOT_CONFLICT" {
		t.Fatalf("expected SLOT_CONFLICT, got %s", apiErr.Code)
	}
}

func TestCorsPreflightReturns204AndHeaders(t *testing.T) {
	now := time.Date(2026, 3, 31, 6, 0, 0, 0, time.UTC)
	svc := service.NewCalendarService(func() time.Time { return now })
	server := NewServer(svc)
	handler := WithCORS(api.Handler(server))

	req := httptest.NewRequest(http.MethodOptions, "/public/name-owner", nil)
	req.Header.Set("Origin", "http://127.0.0.1:5174")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", recorder.Code)
	}

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://127.0.0.1:5174" {
		t.Fatalf("expected allow origin header, got %q", got)
	}
}

func TestAdminUpdateOwnerProfile(t *testing.T) {
	now := time.Date(2026, 3, 31, 6, 0, 0, 0, time.UTC)
	svc := service.NewCalendarService(func() time.Time { return now })
	server := NewServer(svc)
	handler := api.Handler(server)

	payload := map[string]any{
		"photoUrl":    "https://example.com/new-avatar.jpg",
		"displayName": "Alex Updated",
		"email":       "alex.updated@example.com",
		"timezone":    "UTC",
		"weeklySchedule": []map[string]any{
			{"dayOfWeek": "monday", "isActive": true, "startHour": 10, "endHour": 16},
			{"dayOfWeek": "tuesday", "isActive": true, "startHour": 9, "endHour": 18},
			{"dayOfWeek": "wednesday", "isActive": true, "startHour": 9, "endHour": 18},
			{"dayOfWeek": "thursday", "isActive": true, "startHour": 9, "endHour": 18},
			{"dayOfWeek": "friday", "isActive": true, "startHour": 9, "endHour": 18},
			{"dayOfWeek": "saturday", "isActive": false, "startHour": 9, "endHour": 18},
			{"dayOfWeek": "sunday", "isActive": false, "startHour": 9, "endHour": 18},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/admin/owner", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var owner api.CalendarOwner
	if err := json.Unmarshal(recorder.Body.Bytes(), &owner); err != nil {
		t.Fatalf("failed to unmarshal owner: %v", err)
	}

	if owner.DisplayName != "Alex Updated" {
		t.Fatalf("expected updated display name, got %s", owner.DisplayName)
	}
	if owner.Email != "alex.updated@example.com" {
		t.Fatalf("expected updated email, got %s", owner.Email)
	}
	if owner.PhotoUrl == nil || *owner.PhotoUrl != "https://example.com/new-avatar.jpg" {
		t.Fatalf("expected updated photoUrl")
	}
	if len(owner.WeeklySchedule) != 7 {
		t.Fatalf("expected weekly schedule with 7 days, got %d", len(owner.WeeklySchedule))
	}
}

func TestCorsSimpleRequestAddsAllowOriginHeader(t *testing.T) {
	now := time.Date(2026, 3, 31, 6, 0, 0, 0, time.UTC)
	svc := service.NewCalendarService(func() time.Time { return now })
	server := NewServer(svc)
	handler := WithCORS(api.Handler(server))

	req := httptest.NewRequest(http.MethodGet, "/public/name-owner", nil)
	req.Header.Set("Origin", "http://127.0.0.1:5174")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://127.0.0.1:5174" {
		t.Fatalf("expected allow origin header, got %q", got)
	}
}
