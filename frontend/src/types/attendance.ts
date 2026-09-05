export enum AttendanceStatus { PRESENT, ABSENT, LATE }
export interface AttendanceRecord { id: string; employeeId: string; date: string; status: AttendanceStatus; }
export interface Shift { id: string; name: string; startTime: string; endTime: string; }
export interface ShiftAssignment { id: string; employeeId: string; shiftId: string; date: string; }
export interface ShiftSwapRequest { id: string; }
export interface AttendanceFilters { date?: string; }
