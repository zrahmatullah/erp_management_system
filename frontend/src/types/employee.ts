export enum EmployeeStatus { ACTIVE, INACTIVE, ON_LEAVE }
export interface Employee { id: string; firstName: string; lastName: string; status: EmployeeStatus; }
export interface Department { id: string; name: string; }
export interface Position { id: string; name: string; }
export interface EmploymentHistory { id: string; employeeId: string; }
export interface EmployeeCreateRequest { firstName: string; lastName: string; }
export interface EmployeeFilters { status?: EmployeeStatus; }
