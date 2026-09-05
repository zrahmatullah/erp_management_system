export enum LeaveStatus { PENDING, APPROVED, REJECTED }
export interface LeaveType { id: string; name: string; }
export interface LeaveBalance { id: string; typeId: string; balance: number; }
export interface LeaveRequest { id: string; status: LeaveStatus; }
export interface LeaveRequestCreate { typeId: string; startDate: string; endDate: string; }
export interface LeaveFilters { status?: LeaveStatus; }
