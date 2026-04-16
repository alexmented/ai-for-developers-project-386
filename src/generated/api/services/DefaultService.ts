/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Booking } from '../models/Booking';
import type { CalendarOwner } from '../models/CalendarOwner';
import type { CreateBookingRequest } from '../models/CreateBookingRequest';
import type { CreateEventTypeRequest } from '../models/CreateEventTypeRequest';
import type { EventType } from '../models/EventType';
import type { OwnerProfileUpdateRequest } from '../models/OwnerProfileUpdateRequest';
import type { PublicOwnerProfile } from '../models/PublicOwnerProfile';
import type { Slot } from '../models/Slot';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class DefaultService {
    /**
     * @param from
     * @param to
     * @returns Booking The request has succeeded.
     * @throws ApiError
     */
    public static adminApiListUpcomingBookings(
        from?: string,
        to?: string,
    ): CancelablePromise<Array<Booking>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/admin/bookings/upcoming',
            query: {
                'from': from,
                'to': to,
            },
        });
    }
    /**
     * @param requestBody
     * @returns EventType The request has succeeded.
     * @throws ApiError
     */
    public static adminApiCreateEventType(
        requestBody: CreateEventTypeRequest,
    ): CancelablePromise<EventType> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/admin/event-types',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `The server could not understand the request due to invalid syntax.`,
            },
        });
    }
    /**
     * @returns EventType The request has succeeded.
     * @throws ApiError
     */
    public static adminApiListEventTypes(): CancelablePromise<Array<EventType>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/admin/event-types',
        });
    }
    /**
     * @returns CalendarOwner The request has succeeded.
     * @throws ApiError
     */
    public static adminApiGetOwner(): CancelablePromise<CalendarOwner> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/admin/owner',
        });
    }
    /**
     * @param requestBody
     * @returns CalendarOwner The request has succeeded.
     * @throws ApiError
     */
    public static adminApiUpdateOwner(
        requestBody: OwnerProfileUpdateRequest,
    ): CancelablePromise<CalendarOwner> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/admin/owner',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `The server could not understand the request due to invalid syntax.`,
            },
        });
    }
    /**
     * @param ownerSlug
     * @returns PublicOwnerProfile The request has succeeded.
     * @throws ApiError
     */
    public static publicApiGetOwnerProfile(
        ownerSlug: string,
    ): CancelablePromise<PublicOwnerProfile> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/public/{ownerSlug}',
            path: {
                'ownerSlug': ownerSlug,
            },
            errors: {
                404: `The server cannot find the requested resource.`,
            },
        });
    }
    /**
     * @param ownerSlug
     * @param requestBody
     * @returns Booking The request has succeeded.
     * @throws ApiError
     */
    public static publicApiCreateBooking(
        ownerSlug: string,
        requestBody: CreateBookingRequest,
    ): CancelablePromise<Booking> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/public/{ownerSlug}/bookings',
            path: {
                'ownerSlug': ownerSlug,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `The server could not understand the request due to invalid syntax.`,
                404: `The server cannot find the requested resource.`,
                409: `The request conflicts with the current state of the server.`,
            },
        });
    }
    /**
     * @param ownerSlug
     * @returns EventType The request has succeeded.
     * @throws ApiError
     */
    public static publicApiListPublicEventTypes(
        ownerSlug: string,
    ): CancelablePromise<Array<EventType>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/public/{ownerSlug}/event-types',
            path: {
                'ownerSlug': ownerSlug,
            },
            errors: {
                404: `The server cannot find the requested resource.`,
            },
        });
    }
    /**
     * @param ownerSlug
     * @param eventTypeId
     * @param from
     * @param to
     * @returns Slot The request has succeeded.
     * @throws ApiError
     */
    public static publicApiListAvailableSlots(
        ownerSlug: string,
        eventTypeId: string,
        from?: string,
        to?: string,
    ): CancelablePromise<Array<Slot>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/public/{ownerSlug}/event-types/{eventTypeId}/slots',
            path: {
                'ownerSlug': ownerSlug,
                'eventTypeId': eventTypeId,
            },
            query: {
                'from': from,
                'to': to,
            },
            errors: {
                404: `The server cannot find the requested resource.`,
            },
        });
    }
}
