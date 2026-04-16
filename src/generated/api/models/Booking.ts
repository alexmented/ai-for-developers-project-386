/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BookingStatus } from './BookingStatus';
export type Booking = {
    id: string;
    ownerSlug: string;
    eventTypeId: string;
    startAt: string;
    endAt: string;
    guestName: string;
    guestEmail: string;
    guestComment?: string;
    status: BookingStatus;
    createdAt: string;
};

