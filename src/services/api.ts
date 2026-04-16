import type {
  Booking,
  CalendarOwner,
  CreateBookingRequest,
  CreateEventTypeRequest,
  EventType,
  OwnerProfileUpdateRequest,
  PublicOwnerProfile,
  Slot,
} from '@/types/api'
import { normalizeApiError } from '@/types/api'
import { DefaultService, OpenAPI } from '@/generated/api'

OpenAPI.BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:4020'

export const api = {
  async getAdminOwner(): Promise<CalendarOwner> {
    try {
      return await DefaultService.adminApiGetOwner()
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async updateAdminOwner(payload: OwnerProfileUpdateRequest): Promise<CalendarOwner> {
    try {
      return await DefaultService.adminApiUpdateOwner(payload)
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async getOwnerProfile(ownerSlug: string): Promise<PublicOwnerProfile> {
    try {
      return await DefaultService.publicApiGetOwnerProfile(ownerSlug)
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async getPublicEventTypes(ownerSlug: string): Promise<EventType[]> {
    try {
      return await DefaultService.publicApiListPublicEventTypes(ownerSlug)
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async getAvailableSlots(ownerSlug: string, eventTypeId: string, from: string, to: string): Promise<Slot[]> {
    try {
      return await DefaultService.publicApiListAvailableSlots(ownerSlug, eventTypeId, from, to)
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async createBooking(ownerSlug: string, payload: CreateBookingRequest): Promise<Booking> {
    try {
      return await DefaultService.publicApiCreateBooking(ownerSlug, payload)
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async getAdminEventTypes(): Promise<EventType[]> {
    try {
      return await DefaultService.adminApiListEventTypes()
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async createAdminEventType(payload: CreateEventTypeRequest): Promise<EventType> {
    try {
      return await DefaultService.adminApiCreateEventType(payload)
    } catch (error) {
      throw normalizeApiError(error)
    }
  },

  async getAdminUpcomingBookings(): Promise<Booking[]> {
    try {
      return await DefaultService.adminApiListUpcomingBookings()
    } catch (error) {
      throw normalizeApiError(error)
    }
  },
}
