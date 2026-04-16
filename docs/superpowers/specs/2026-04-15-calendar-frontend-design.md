---
title: Calendar Frontend Design (Vite + Vue + shadcn-vue)
date: 2026-04-15
status: approved-by-user-in-chat
---

## 1. Goal and Scope

Implement the frontend for a booking application based on provided UI references, preserving visual style and interaction patterns.

In scope:
- Guest landing page.
- Guest event-type selection page.
- Guest slot selection and booking page.
- Admin page in the same visual style.
- Responsive behavior (desktop/tablet/mobile) and pixel-accurate spacing/visual hierarchy.
- Integration with mock API via Prism.

Out of scope:
- Real backend integration.
- Authentication and authorization flows.

## 2. Stack and Constraints

- Build tool: Vite 8.
- Framework: Vue 3.
- Language: TypeScript.
- UI primitives: shadcn-vue components.
- Routing: vue-router.
- Data source: Prism mock server generated from API contract.

Design constraints:
- Keep the same visual language as provided mockups.
- Maintain CTA emphasis and card-first composition.
- Keep interactions minimal and predictable.

## 3. Information Architecture and Routes

Application routes:
- `/` — landing page.
- `/:ownerSlug` — public owner profile and event type selection.
- `/:ownerSlug/:eventTypeId` — date/slot selection and booking action.
- `/admin` — event type management and upcoming bookings.

Navigation:
- Header includes links to guest flow (`Записаться`) and admin (`Админка`).
- Guest flow preserves `ownerSlug` in route context.

## 4. API Integration Contract (Prism)

Frontend uses these endpoints:

Public:
- `GET /public/{ownerSlug}`
- `GET /public/{ownerSlug}/event-types`
- `GET /public/{ownerSlug}/event-types/{eventTypeId}/slots?from&to`
- `POST /public/{ownerSlug}/bookings`

Admin:
- `GET /admin/owner`
- `GET /admin/event-types`
- `POST /admin/event-types`
- `GET /admin/bookings/upcoming`

Business constraints reflected in UI:
- Slots are queried in the booking window of 14 days from current date by default.
- Booking conflicts (`409 SLOT_CONFLICT`) are handled with user feedback and slot refresh.

## 5. UI Composition

### 5.1 Global Layout

- Top navigation bar with logo and two links.
- Soft neutral background and lightly elevated cards.
- Rounded corners, subtle borders, compact badges.

### 5.2 Landing (`/`)

- Left hero block:
  - label badge,
  - title,
  - short supporting text,
  - primary CTA button.
- Right capabilities card with concise bullet list.

### 5.3 Event Types (`/:ownerSlug`)

- Intro card with host avatar/name/role.
- Main title and short helper text.
- Grid/list of event-type cards:
  - title,
  - description,
  - duration badge.
- Card click routes to slot selection page.

### 5.4 Slot Selection (`/:ownerSlug/:eventTypeId`)

Three-column desktop layout:
1) Selected event card and summary fields (chosen date/time).
2) Calendar grid (month navigation and day cells).
3) Slot status panel with actionable list and footer buttons.

Interaction:
- Selecting a day loads or filters slots.
- Selecting a slot updates summary.
- Continue triggers booking request.

### 5.5 Admin (`/admin`)

- Event type creation form:
  - id,
  - name,
  - description,
  - duration in minutes.
- Event type list for quick verification.
- Upcoming bookings list/table across all event types.

## 6. Component Plan

Reusable components:
- `AppHeader`
- `HeroSection`
- `OwnerIntroCard`
- `EventTypeCard`
- `BookingSummaryCard`
- `MiniCalendarGrid`
- `SlotStatusList`
- `AdminEventTypeForm`
- `AdminUpcomingBookingsList`

shadcn-vue base components (initial set):
- Button
- Card
- Badge
- Input
- Textarea
- Table
- Dialog (optional for confirmations)

## 7. State and Data Flow

- A lightweight API service layer centralizes fetch calls and DTO mapping.
- Pages manage local screen state (`loading`, `error`, `empty`, `ready`).
- Route params (`ownerSlug`, `eventTypeId`) are source of truth for current context.
- Booking payload uses selected slot start-time and guest fields.

Default query strategy:
- `from = now`
- `to = now + 14 days`

## 8. Error Handling and UX Rules

- `OWNER_NOT_FOUND` → dedicated not-found state on public pages.
- `NOT_FOUND` for event type → back link to owner event list.
- `VALIDATION_ERROR` → field-level messages on admin and booking forms.
- `SLOT_CONFLICT` → inline/toast message + slot list refresh.

Empty states:
- No event types.
- No available slots for selected date.
- No upcoming bookings in admin.

## 9. Responsive and Pixel-Accuracy Strategy

Breakpoints:
- Desktop: full composition as in references.
- Tablet: 2-column adaptation, slots panel moves below calendar if needed.
- Mobile: single-column stacked flow, preserved hierarchy and CTA prominence.

Pixel-accuracy checklist:
- Typography scale per screen role.
- Card radius/border/shadow consistency.
- Spacing rhythm across sections.
- Controls height and alignment.
- Badge and button optical alignment.

## 10. Verification Plan

- Run type/build checks after implementation.
- Manual route walkthrough:
  - `/`
  - `/:ownerSlug`
  - `/:ownerSlug/:eventTypeId`
  - `/admin`
- Validate responsive behavior on representative widths.
- Validate conflict and not-found API scenarios.

## 11. Implementation Readiness

This design is intentionally implementation-ready with fixed scope, route map, UI composition, component boundaries, API touchpoints, and validation criteria. Next step is generating an execution plan and then implementing screens/components in the agreed stack.
