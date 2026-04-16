package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexmented/ai-for-developers-project-386/backend/internal/api"
	"github.com/alexmented/ai-for-developers-project-386/backend/internal/service"
)

type Server struct {
	svc *service.CalendarService
}

func NewServer(svc *service.CalendarService) *Server {
	return &Server{svc: svc}
}

func (s *Server) AdminApiListUpcomingBookings(w http.ResponseWriter, _ *http.Request, params api.AdminApiListUpcomingBookingsParams) {
	bookings, err := s.svc.ListUpcomingBookings(params.From, params.To)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, bookings)
}

func (s *Server) AdminApiListEventTypes(w http.ResponseWriter, _ *http.Request) {
	eventTypes, err := s.svc.ListEventTypes("name-owner")
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, eventTypes)
}

func (s *Server) AdminApiCreateEventType(w http.ResponseWriter, r *http.Request) {
	var req api.CreateEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid json body")
		return
	}

	eventType, err := s.svc.CreateEventType(service.CreateEventTypeInput{
		ID:              req.Id,
		Name:            req.Name,
		Description:     stringValue(req.Description),
		DurationMinutes: req.DurationMinutes,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, eventType)
}

func (s *Server) AdminApiGetOwner(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.svc.GetOwner())
}

func (s *Server) AdminApiUpdateOwner(w http.ResponseWriter, r *http.Request) {
	var req api.OwnerProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid json body")
		return
	}

	owner, err := s.svc.UpdateOwnerProfile(service.UpdateOwnerProfileInput{
		PhotoURL:       req.PhotoUrl,
		DisplayName:    req.DisplayName,
		Email:          req.Email,
		Timezone:       req.Timezone,
		WeeklySchedule: req.WeeklySchedule,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, owner)
}

func (s *Server) PublicApiGetOwnerProfile(w http.ResponseWriter, _ *http.Request, ownerSlug string) {
	profile, err := s.svc.GetOwnerProfile(ownerSlug)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

func (s *Server) PublicApiCreateBooking(w http.ResponseWriter, r *http.Request, ownerSlug string) {
	var req api.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid json body")
		return
	}

	booking, err := s.svc.CreateBooking(ownerSlug, service.CreateBookingInput{
		EventTypeID:  req.EventTypeId,
		SlotStartAt:  req.SlotStartAt,
		GuestName:    req.GuestName,
		GuestEmail:   req.GuestEmail,
		GuestComment: req.GuestComment,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, booking)
}

func (s *Server) PublicApiListPublicEventTypes(w http.ResponseWriter, _ *http.Request, ownerSlug string) {
	eventTypes, err := s.svc.ListEventTypes(ownerSlug)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, eventTypes)
}

func (s *Server) PublicApiListAvailableSlots(w http.ResponseWriter, _ *http.Request, ownerSlug string, eventTypeId string, params api.PublicApiListAvailableSlotsParams) {
	slots, err := s.svc.ListAvailableSlots(ownerSlug, eventTypeId, params.From, params.To)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, slots)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, api.ErrorResponse{Code: code, Message: message})
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func writeServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {
		writeError(w, svcErr.Status, svcErr.Code, svcErr.Message)
		return
	}

	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func NewDefaultHandler() http.Handler {
	now := func() time.Time { return time.Now().UTC() }
	svc := service.NewCalendarService(now)
	server := NewServer(svc)
	return WithCORS(api.Handler(server))
}

func WithCORS(next http.Handler) http.Handler {
	allowedMethods := "GET,POST,PUT,OPTIONS"
	allowedHeaders := "Content-Type,Authorization"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
