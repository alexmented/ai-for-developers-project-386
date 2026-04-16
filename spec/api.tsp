import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

@service(#{
  title: "Calendar Booking API"
})
namespace CalendarBooking;

alias OwnerSlug = string;
alias EventTypeId = string;
alias BookingId = string;

enum BookingStatus {
  confirmed
}

enum DayOfWeek {
  monday,
  tuesday,
  wednesday,
  thursday,
  friday,
  saturday,
  sunday
}

model WorkDaySchedule {
  dayOfWeek: DayOfWeek;
  isActive: boolean;
  startHour: int32;
  endHour: int32;
}

model CalendarOwner {
  id: "owner-default";
  slug: OwnerSlug;
  photoUrl?: string;
  displayName: string;
  email: string;
  timezone: string;
  weeklySchedule: WorkDaySchedule[];
}

model PublicOwnerProfile {
  slug: OwnerSlug;
  photoUrl?: string;
  displayName: string;
  timezone: string;
  weeklySchedule: WorkDaySchedule[];
}

model OwnerProfileUpdateRequest {
  photoUrl?: string;
  displayName: string;
  email: string;
  timezone: string;
  weeklySchedule: WorkDaySchedule[];
}

model EventType {
  id: EventTypeId;
  name: string;
  description: string;
  durationMinutes: int32;
}

model Slot {
  eventTypeId: EventTypeId;
  startAt: utcDateTime;
  endAt: utcDateTime;
  isAvailable: true;
}

model Booking {
  id: BookingId;
  ownerSlug: OwnerSlug;
  eventTypeId: EventTypeId;
  startAt: utcDateTime;
  endAt: utcDateTime;
  guestName: string;
  guestEmail: string;
  guestComment?: string;
  status: BookingStatus;
  createdAt: utcDateTime;
}

model CreateEventTypeRequest {
  id?: EventTypeId;
  name: string;
  description?: string;
  durationMinutes: int32;
}

model CreateBookingRequest {
  eventTypeId: EventTypeId;
  slotStartAt: utcDateTime;
  guestName: string;
  guestEmail: string;
  guestComment?: string;
}

model ErrorResponse {
  code: string;
  message: string;
}

@error
model ValidationError extends ErrorResponse {
  @statusCode statusCode: 400;
  code: "VALIDATION_ERROR";
}

@error
model NotFoundError extends ErrorResponse {
  @statusCode statusCode: 404;
  code: "NOT_FOUND";
}

@error
model OwnerNotFoundError extends ErrorResponse {
  @statusCode statusCode: 404;
  code: "OWNER_NOT_FOUND";
}

@error
model SlotConflictError extends ErrorResponse {
  @statusCode statusCode: 409;
  code: "SLOT_CONFLICT";
}

@route("/admin")
interface AdminApi {
  @get
  @route("/owner")
  op getOwner(): CalendarOwner;

  @put
  @route("/owner")
  op updateOwner(@body input: OwnerProfileUpdateRequest): CalendarOwner | ValidationError;

  @post
  @route("/event-types")
  op createEventType(@body input: CreateEventTypeRequest): EventType | ValidationError;

  @get
  @route("/event-types")
  op listEventTypes(): EventType[];

  @get
  @route("/bookings/upcoming")
  op listUpcomingBookings(
    @query from?: utcDateTime,
    @query to?: utcDateTime
  ): Booking[];
}

@route("/public/{ownerSlug}")
interface PublicApi {
  @get
  op getOwnerProfile(@path ownerSlug: OwnerSlug): PublicOwnerProfile | OwnerNotFoundError;

  @get
  @route("/event-types")
  op listPublicEventTypes(@path ownerSlug: OwnerSlug): EventType[] | OwnerNotFoundError;

  @get
  @route("/event-types/{eventTypeId}/slots")
  op listAvailableSlots(
    @path ownerSlug: OwnerSlug,
    @path eventTypeId: EventTypeId,
    @query from?: utcDateTime,
    @query to?: utcDateTime
  ): Slot[] | OwnerNotFoundError | NotFoundError;

  @post
  @route("/bookings")
  op createBooking(
    @path ownerSlug: OwnerSlug,
    @body input: CreateBookingRequest
  ): Booking | ValidationError | OwnerNotFoundError | NotFoundError | SlotConflictError;
}
