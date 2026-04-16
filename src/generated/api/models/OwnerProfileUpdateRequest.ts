/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { WorkDaySchedule } from './WorkDaySchedule';
export type OwnerProfileUpdateRequest = {
    photoUrl?: string;
    displayName: string;
    email: string;
    timezone: string;
    weeklySchedule: Array<WorkDaySchedule>;
};

