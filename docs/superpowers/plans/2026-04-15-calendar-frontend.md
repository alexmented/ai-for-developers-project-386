# Calendar Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a pixel-perfect responsive frontend (guest + admin) on Vue 3 + Vite 8 + shadcn-vue with Prism-backed API mocks.

**Architecture:** Route-first SPA with four screens (`/`, `/:ownerSlug`, `/:ownerSlug/:eventTypeId`, `/admin`). API calls are centralized in `src/services/api.ts`, while each page handles loading/error/empty states locally. Visual consistency is enforced through shared tokens and reusable UI blocks.

**Tech Stack:** Vite 8, Vue 3, TypeScript, Vue Router, Tailwind CSS, shadcn-vue, Vitest, Vue Test Utils, Prism

---

## File Structure

- Create: `package.json`, `tsconfig.json`, `vite.config.ts`, `index.html`
- Create: `src/main.ts`, `src/App.vue`, `src/router/index.ts`
- Create: `src/types/api.ts`, `src/services/api.ts`, `src/lib/date.ts`
- Create: `src/components/layout/AppHeader.vue`
- Create: `src/components/booking/EventTypeCard.vue`, `src/components/booking/MiniCalendarGrid.vue`, `src/components/booking/SlotStatusList.vue`
- Create: `src/pages/LandingPage.vue`, `src/pages/PublicEventTypesPage.vue`, `src/pages/PublicBookingPage.vue`, `src/pages/AdminPage.vue`
- Create: `src/styles/globals.css`
- Create: `tests/date-window.spec.ts`, `tests/api-mapping.spec.ts`
- Create: `components.json` (shadcn-vue config)

---

### Task 1: Bootstrap Vue + Vite + TS

**Files:** create core project/config files.

- [ ] Step 1: Initialize Vite Vue TS project skeleton.
- [ ] Step 2: Add scripts in `package.json`:

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc -b && vite build",
    "preview": "vite preview",
    "test": "vitest run"
  }
}
```

- [ ] Step 3: Verify bootstrap.

Run: `npm run build`
Expected: successful build with generated `dist`.

- [ ] Step 4: Commit.

```bash
git add package.json tsconfig.json vite.config.ts src/main.ts src/App.vue
git commit -m "chore: bootstrap vite vue typescript app"
```

### Task 2: Configure Tailwind + shadcn-vue base

**Files:** `components.json`, Tailwind setup, base styles.

- [ ] Step 1: Configure Tailwind and shared CSS variables in `src/styles/globals.css`.
- [ ] Step 2: Add shadcn-vue registry configuration and install base components: `button`, `card`, `badge`, `input`, `textarea`, `table`.
- [ ] Step 3: Wire global styles in `src/main.ts`.
- [ ] Step 4: Verify UI primitives render on a temporary smoke component.
- [ ] Step 5: Commit.

### Task 3: Routing and layout shell

**Files:** `src/router/index.ts`, `src/App.vue`, `src/components/layout/AppHeader.vue`.

- [ ] Step 1: Create route table for `/`, `/:ownerSlug`, `/:ownerSlug/:eventTypeId`, `/admin`.
- [ ] Step 2: Implement header with links `Записаться` and `Админка`.
- [ ] Step 3: Add responsive container and global page background.
- [ ] Step 4: Run app and verify navigation.

Run: `npm run dev`
Expected: all routes open without runtime errors.

- [ ] Step 5: Commit.

### Task 4: API service + domain helpers (TDD)

**Files:** `src/types/api.ts`, `src/services/api.ts`, `src/lib/date.ts`, `tests/date-window.spec.ts`, `tests/api-mapping.spec.ts`.

- [ ] Step 1: Write failing test for 14-day window helper.

```ts
import { describe, it, expect } from "vitest";
import { buildBookingWindow } from "../src/lib/date";

describe("buildBookingWindow", () => {
  it("returns now and now+14 days", () => {
    const now = new Date("2026-04-15T10:00:00.000Z");
    const { from, to } = buildBookingWindow(now);
    expect(from).toBe("2026-04-15T10:00:00.000Z");
    expect(to).toBe("2026-04-29T10:00:00.000Z");
  });
});
```

- [ ] Step 2: Run `npm test` and confirm fail.
- [ ] Step 3: Implement `buildBookingWindow` and typed API methods.
- [ ] Step 4: Run `npm test` and confirm pass.
- [ ] Step 5: Commit.

### Task 5: Implement guest pages (pixel-perfect + responsive)

**Files:** guest pages/components.

- [ ] Step 1: Build `LandingPage.vue` to match hero + capabilities card layout.
- [ ] Step 2: Build `PublicEventTypesPage.vue` with owner card and event-type cards.
- [ ] Step 3: Build `PublicBookingPage.vue` with 3-column desktop layout and responsive collapse.
- [ ] Step 4: Connect API calls and states (`loading/error/empty`).
- [ ] Step 5: Handle `409 SLOT_CONFLICT` with message + slot reload.
- [ ] Step 6: Commit.

### Task 6: Implement admin page

**Files:** `src/pages/AdminPage.vue` and helper blocks.

- [ ] Step 1: Add event type creation form with validation messages.
- [ ] Step 2: Add event types list.
- [ ] Step 3: Add upcoming bookings table.
- [ ] Step 4: Bind create/list actions to API service.
- [ ] Step 5: Commit.

### Task 7: Verification and polish

**Files:** styles + small fixes.

- [ ] Step 1: Run pixel pass (spacing, typography, card radii, button heights).
- [ ] Step 2: Verify mobile/tablet breakpoints manually.
- [ ] Step 3: Run checks.

Run:
- `npm run test`
- `npm run build`

Expected: all pass.

- [ ] Step 4: Commit final polish.

```bash
git add .
git commit -m "feat: implement calendar frontend guest and admin flows"
```

## Self-review

- Spec coverage: landing, event types, booking screen, admin, Prism integration, 14-day window, conflict handling, responsive behavior — all covered.
- Placeholder scan: no TODO/TBD placeholders in tasks.
- Type consistency: `ownerSlug`, `eventTypeId`, booking window and error codes are consistent with agreed API contract.

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-04-15-calendar-frontend.md`. Two execution options:

1. Subagent-Driven (recommended) - I dispatch a fresh subagent per task, review between tasks, fast iteration
2. Inline Execution - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?
