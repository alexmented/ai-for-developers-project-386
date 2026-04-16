import { ApiError as GeneratedApiError } from '@/generated/api'
import type { ErrorResponse } from '@/generated/api'

export type {
  Booking,
  CalendarOwner,
  CreateBookingRequest,
  CreateEventTypeRequest,
  DayOfWeek,
  EventType,
  OwnerProfileUpdateRequest,
  PublicOwnerProfile,
  Slot,
  WorkDaySchedule,
} from '@/generated/api'

export class ApiError extends Error {
  code: string
  status: number

  constructor(message: string, code: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

export function normalizeApiError(error: unknown): ApiError {
  if (error instanceof ApiError) {
    return error
  }

  if (error instanceof GeneratedApiError) {
    const body = (error.body ?? {}) as Partial<ErrorResponse>
    return new ApiError(body.message ?? error.message ?? 'Request failed', body.code ?? 'UNKNOWN_ERROR', error.status)
  }

  if (error instanceof Error) {
    return new ApiError(error.message, 'UNKNOWN_ERROR', 500)
  }

  return new ApiError('Request failed', 'UNKNOWN_ERROR', 500)
}
