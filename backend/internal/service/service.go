package service

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/alexmented/ai-for-developers-project-386/backend/internal/api"
)

const (
	ownerID       = api.OwnerDefault
	ownerSlug     = "name-owner"
	ownerName     = "Tota"
	ownerEmail    = "owner@example.com"
	ownerPhotoURL = "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?auto=format&fit=crop&w=400&q=80"
	ownerTimezone = "Europe/Moscow"

	bookingWindowDays = 14
	slotStepMinutes   = 15
	defaultWorkStart  = 9
	defaultWorkEnd    = 18
)

type nowFunc func() time.Time

type ServiceError struct {
	Code    string
	Message string
	Status  int
}

func (e *ServiceError) Error() string {
	return e.Message
}

type CreateBookingInput struct {
	EventTypeID  string
	SlotStartAt  time.Time
	GuestName    string
	GuestEmail   string
	GuestComment *string
}

type CreateEventTypeInput struct {
	ID              *string
	Name            string
	Description     string
	DurationMinutes int32
}

type UpdateOwnerProfileInput struct {
	PhotoURL       *string
	DisplayName    string
	Email          string
	Timezone       string
	WeeklySchedule []api.WorkDaySchedule
}

type CalendarService struct {
	now nowFunc

	mu            sync.RWMutex
	owner         api.CalendarOwner
	location      *time.Location
	eventTypes    map[string]api.EventType
	bookings      []api.Booking
	nextBookingID int
}

func NewCalendarService(now nowFunc) *CalendarService {
	if now == nil {
		now = time.Now
	}

	location, err := time.LoadLocation(ownerTimezone)
	if err != nil {
		location = time.UTC
	}

	eventTypes := map[string]api.EventType{
		"meeting-15": {
			Id:              "meeting-15",
			Name:            "Встреча 15 минут",
			Description:     "Короткий тип события для быстрого слота.",
			DurationMinutes: 15,
		},
		"meeting-30": {
			Id:              "meeting-30",
			Name:            "Встреча 30 минут",
			Description:     "Базовый тип события для бронирования.",
			DurationMinutes: 30,
		},
	}

	return &CalendarService{
		now: now,
		owner: api.CalendarOwner{
			Id:             ownerID,
			Slug:           ownerSlug,
			PhotoUrl:       stringPtr(ownerPhotoURL),
			DisplayName:    ownerName,
			Email:          ownerEmail,
			Timezone:       ownerTimezone,
			WeeklySchedule: defaultWeeklySchedule(),
		},
		location:      location,
		eventTypes:    eventTypes,
		bookings:      make([]api.Booking, 0),
		nextBookingID: 1,
	}
}

func (s *CalendarService) GetOwner() api.CalendarOwner {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.owner
}

func (s *CalendarService) GetOwnerProfile(slug string) (api.PublicOwnerProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if slug != s.owner.Slug {
		return api.PublicOwnerProfile{}, ownerNotFoundError(slug)
	}

	return api.PublicOwnerProfile{
		Slug:           s.owner.Slug,
		PhotoUrl:       s.owner.PhotoUrl,
		DisplayName:    s.owner.DisplayName,
		Timezone:       s.owner.Timezone,
		WeeklySchedule: cloneWeeklySchedule(s.owner.WeeklySchedule),
	}, nil
}

func (s *CalendarService) UpdateOwnerProfile(input UpdateOwnerProfileInput) (api.CalendarOwner, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	photoURL := (*string)(nil)
	if input.PhotoURL != nil {
		trimmed := strings.TrimSpace(*input.PhotoURL)
		if trimmed != "" {
			parsed, err := url.ParseRequestURI(trimmed)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" {
				return api.CalendarOwner{}, validationError("photoUrl должен быть корректным URL")
			}
			photoURL = &trimmed
		}
	}

	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		return api.CalendarOwner{}, validationError("displayName обязателен")
	}

	email := strings.TrimSpace(input.Email)
	if email == "" {
		return api.CalendarOwner{}, validationError("email обязателен")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return api.CalendarOwner{}, validationError("email должен быть корректным")
	}

	timezone := strings.TrimSpace(input.Timezone)
	if timezone == "" {
		return api.CalendarOwner{}, validationError("timezone обязателен")
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return api.CalendarOwner{}, validationError("timezone должен быть валидным IANA идентификатором")
	}

	weeklySchedule, err := normalizeWeeklySchedule(input.WeeklySchedule)
	if err != nil {
		return api.CalendarOwner{}, err
	}

	s.owner.PhotoUrl = photoURL
	s.owner.DisplayName = displayName
	s.owner.Email = email
	s.owner.Timezone = timezone
	s.owner.WeeklySchedule = weeklySchedule
	s.location = location

	return s.owner, nil
}

func (s *CalendarService) ListEventTypes(slug string) ([]api.EventType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if slug != s.owner.Slug {
		return nil, ownerNotFoundError(slug)
	}

	items := make([]api.EventType, 0, len(s.eventTypes))
	for _, item := range s.eventTypes {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Id < items[j].Id
	})

	return items, nil
}

func (s *CalendarService) CreateEventType(input CreateEventTypeInput) (api.EventType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)
	if name == "" || input.DurationMinutes <= 0 {
		return api.EventType{}, validationError("name и durationMinutes обязательны")
	}

	eventTypeID := ""
	if input.ID != nil {
		eventTypeID = strings.TrimSpace(*input.ID)
	}
	if eventTypeID == "" {
		eventTypeID = s.generateEventTypeID(name)
	}

	if _, exists := s.eventTypes[eventTypeID]; exists {
		return api.EventType{}, validationError("тип события с таким id уже существует")
	}

	eventType := api.EventType{
		Id:              eventTypeID,
		Name:            name,
		Description:     description,
		DurationMinutes: input.DurationMinutes,
	}

	s.eventTypes[eventTypeID] = eventType

	return eventType, nil
}

func (s *CalendarService) ListUpcomingBookings(from, to *time.Time) ([]api.Booking, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	windowStart, windowEnd, err := s.resolveWindow(from, to)
	if err != nil {
		return nil, err
	}

	items := make([]api.Booking, 0)
	for _, booking := range s.bookings {
		if booking.StartAt.Before(windowStart) {
			continue
		}
		if booking.StartAt.After(windowEnd) {
			continue
		}
		items = append(items, booking)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].StartAt.Before(items[j].StartAt)
	})

	return items, nil
}

func (s *CalendarService) ListAvailableSlots(slug, eventTypeID string, from, to *time.Time) ([]api.Slot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if slug != s.owner.Slug {
		return nil, ownerNotFoundError(slug)
	}

	eventType, ok := s.eventTypes[eventTypeID]
	if !ok {
		return nil, notFoundError("тип события не найден")
	}

	windowStart, windowEnd, err := s.resolveWindow(from, to)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(eventType.DurationMinutes) * time.Minute
	step := duration

	startLocal := windowStart.In(s.location)
	endLocal := windowEnd.In(s.location)
	cursorDay := time.Date(startLocal.Year(), startLocal.Month(), startLocal.Day(), 0, 0, 0, 0, s.location)
	lastDay := time.Date(endLocal.Year(), endLocal.Month(), endLocal.Day(), 0, 0, 0, 0, s.location)

	slots := make([]api.Slot, 0)
	for !cursorDay.After(lastDay) {
		daySchedule, ok := scheduleForWeekday(s.owner.WeeklySchedule, cursorDay.Weekday())
		if !ok || !daySchedule.IsActive {
			cursorDay = cursorDay.AddDate(0, 0, 1)
			continue
		}

		dayStart := time.Date(cursorDay.Year(), cursorDay.Month(), cursorDay.Day(), int(daySchedule.StartHour), 0, 0, 0, s.location)
		dayEnd := time.Date(cursorDay.Year(), cursorDay.Month(), cursorDay.Day(), int(daySchedule.EndHour), 0, 0, 0, s.location)

		for slotStart := dayStart; slotStart.Add(duration).Equal(dayEnd) || slotStart.Add(duration).Before(dayEnd); slotStart = slotStart.Add(step) {
			slotStartUTC := slotStart.UTC()
			slotEndUTC := slotStart.Add(duration).UTC()

			if slotEndUTC.Before(windowStart) || slotEndUTC.Equal(windowStart) {
				continue
			}
			if slotStartUTC.After(windowEnd) || slotStartUTC.Equal(windowEnd) {
				continue
			}

			slots = append(slots, api.Slot{
				EventTypeId: eventTypeID,
				StartAt:     slotStartUTC,
				EndAt:       slotEndUTC,
				IsAvailable: api.SlotIsAvailable(!s.hasConflict(slotStartUTC, slotEndUTC, "")),
			})
		}

		cursorDay = cursorDay.AddDate(0, 0, 1)
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].StartAt.Before(slots[j].StartAt)
	})

	return slots, nil
}

func (s *CalendarService) CreateBooking(slug string, input CreateBookingInput) (api.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if slug != s.owner.Slug {
		return api.Booking{}, ownerNotFoundError(slug)
	}

	eventType, ok := s.eventTypes[input.EventTypeID]
	if !ok {
		return api.Booking{}, notFoundError("тип события не найден")
	}

	if strings.TrimSpace(input.GuestName) == "" || strings.TrimSpace(input.GuestEmail) == "" {
		return api.Booking{}, validationError("guestName и guestEmail обязательны")
	}

	windowStart := s.now().UTC()
	windowEnd := windowStart.AddDate(0, 0, bookingWindowDays)
	slotStart := input.SlotStartAt.UTC()
	if slotStart.Before(windowStart) || slotStart.After(windowEnd) {
		return api.Booking{}, validationError("слот должен быть в окне ближайших 14 дней")
	}

	duration := time.Duration(eventType.DurationMinutes) * time.Minute
	slotEnd := slotStart.Add(duration)
	if !s.inWorkingHours(slotStart, slotEnd) {
		return api.Booking{}, validationError("слот вне доступного расписания дня")
	}

	if slotStart.Minute()%int(eventType.DurationMinutes) != 0 || slotStart.Second() != 0 || slotStart.Nanosecond() != 0 {
		return api.Booking{}, validationError(fmt.Sprintf("слот должен быть кратен %d минутам", eventType.DurationMinutes))
	}

	if s.hasConflict(slotStart, slotEnd, "") {
		return api.Booking{}, conflictError("слот уже занят")
	}

	bookingID := fmt.Sprintf("bk-%03d", s.nextBookingID)
	s.nextBookingID++

	booking := api.Booking{
		Id:           bookingID,
		OwnerSlug:    s.owner.Slug,
		EventTypeId:  input.EventTypeID,
		StartAt:      slotStart,
		EndAt:        slotEnd,
		GuestName:    input.GuestName,
		GuestEmail:   input.GuestEmail,
		GuestComment: input.GuestComment,
		Status:       api.Confirmed,
		CreatedAt:    s.now().UTC(),
	}

	s.bookings = append(s.bookings, booking)

	return booking, nil
}

func (s *CalendarService) resolveWindow(from, to *time.Time) (time.Time, time.Time, error) {
	windowStart := s.now().UTC()
	windowEnd := windowStart.AddDate(0, 0, bookingWindowDays)

	if from != nil {
		windowStart = from.UTC()
	}
	if to != nil {
		windowEnd = to.UTC()
	}

	if !windowEnd.After(windowStart) {
		return time.Time{}, time.Time{}, validationError("to должно быть позже from")
	}

	return windowStart, windowEnd, nil
}

func (s *CalendarService) hasConflict(start, end time.Time, bookingID string) bool {
	for _, booking := range s.bookings {
		if bookingID != "" && booking.Id == bookingID {
			continue
		}

		if start.Before(booking.EndAt) && end.After(booking.StartAt) {
			return true
		}
	}

	return false
}

func (s *CalendarService) inWorkingHours(startUTC, endUTC time.Time) bool {
	start := startUTC.In(s.location)
	end := endUTC.In(s.location)

	if start.Year() != end.Year() || start.YearDay() != end.YearDay() {
		return false
	}

	daySchedule, ok := scheduleForWeekday(s.owner.WeeklySchedule, start.Weekday())
	if !ok || !daySchedule.IsActive {
		return false
	}

	dayStart := time.Date(start.Year(), start.Month(), start.Day(), int(daySchedule.StartHour), 0, 0, 0, s.location)
	dayEnd := time.Date(start.Year(), start.Month(), start.Day(), int(daySchedule.EndHour), 0, 0, 0, s.location)

	return (start.Equal(dayStart) || start.After(dayStart)) && (end.Equal(dayEnd) || end.Before(dayEnd))
}

func stringPtr(value string) *string {
	return &value
}

func cloneWeeklySchedule(input []api.WorkDaySchedule) []api.WorkDaySchedule {
	if len(input) == 0 {
		return nil
	}

	output := make([]api.WorkDaySchedule, len(input))
	copy(output, input)
	return output
}

func defaultWeeklySchedule() []api.WorkDaySchedule {
	return []api.WorkDaySchedule{
		{DayOfWeek: api.Monday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
		{DayOfWeek: api.Tuesday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
		{DayOfWeek: api.Wednesday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
		{DayOfWeek: api.Thursday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
		{DayOfWeek: api.Friday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
		{DayOfWeek: api.Saturday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
		{DayOfWeek: api.Sunday, IsActive: true, StartHour: defaultWorkStart, EndHour: defaultWorkEnd},
	}
}

func normalizeWeeklySchedule(input []api.WorkDaySchedule) ([]api.WorkDaySchedule, error) {
	if len(input) != 7 {
		return nil, validationError("weeklySchedule должен содержать 7 дней")
	}

	orderedDays := []api.DayOfWeek{api.Monday, api.Tuesday, api.Wednesday, api.Thursday, api.Friday, api.Saturday, api.Sunday}
	byDay := make(map[api.DayOfWeek]api.WorkDaySchedule, len(input))

	for _, day := range input {
		if !day.DayOfWeek.Valid() {
			return nil, validationError("weeklySchedule содержит неизвестный день недели")
		}
		if _, exists := byDay[day.DayOfWeek]; exists {
			return nil, validationError("weeklySchedule содержит дубли дней недели")
		}

		if day.IsActive {
			if day.StartHour < 0 || day.StartHour > 23 || day.EndHour < 1 || day.EndHour > 24 {
				return nil, validationError("часы активного дня должны быть в диапазоне 00:00-24:00")
			}
			if day.StartHour >= day.EndHour {
				return nil, validationError("endHour должен быть позже startHour для активного дня")
			}
		} else {
			day.StartHour = defaultWorkStart
			day.EndHour = defaultWorkEnd
		}

		byDay[day.DayOfWeek] = day
	}

	result := make([]api.WorkDaySchedule, 0, len(orderedDays))
	for _, day := range orderedDays {
		schedule, exists := byDay[day]
		if !exists {
			return nil, validationError("weeklySchedule должен содержать все дни недели")
		}
		result = append(result, schedule)
	}

	return result, nil
}

func scheduleForWeekday(weeklySchedule []api.WorkDaySchedule, weekday time.Weekday) (api.WorkDaySchedule, bool) {
	targetDay, ok := weekdayToAPIDay(weekday)
	if !ok {
		return api.WorkDaySchedule{}, false
	}

	for _, day := range weeklySchedule {
		if day.DayOfWeek == targetDay {
			return day, true
		}
	}

	return api.WorkDaySchedule{}, false
}

func weekdayToAPIDay(weekday time.Weekday) (api.DayOfWeek, bool) {
	switch weekday {
	case time.Monday:
		return api.Monday, true
	case time.Tuesday:
		return api.Tuesday, true
	case time.Wednesday:
		return api.Wednesday, true
	case time.Thursday:
		return api.Thursday, true
	case time.Friday:
		return api.Friday, true
	case time.Saturday:
		return api.Saturday, true
	case time.Sunday:
		return api.Sunday, true
	default:
		return "", false
	}
}

func (s *CalendarService) generateEventTypeID(name string) string {
	base := slugify(name)
	candidate := base
	index := 2

	for {
		if _, exists := s.eventTypes[candidate]; !exists {
			return candidate
		}

		candidate = fmt.Sprintf("%s-%d", base, index)
		index++
	}
}

func slugify(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "event"
	}

	var builder strings.Builder
	prevDash := false

	for _, char := range trimmed {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			builder.WriteRune(char)
			prevDash = false
			continue
		}

		if !prevDash {
			builder.WriteByte('-')
			prevDash = true
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "event"
	}

	return result
}

func ownerNotFoundError(slug string) error {
	return &ServiceError{Code: "OWNER_NOT_FOUND", Message: fmt.Sprintf("owner %q not found", slug), Status: 404}
}

func notFoundError(message string) error {
	return &ServiceError{Code: "NOT_FOUND", Message: message, Status: 404}
}

func validationError(message string) error {
	return &ServiceError{Code: "VALIDATION_ERROR", Message: message, Status: 400}
}

func conflictError(message string) error {
	return &ServiceError{Code: "SLOT_CONFLICT", Message: message, Status: 409}
}

func AsServiceError(err error) (*ServiceError, bool) {
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		return svcErr, true
	}

	return nil, false
}
