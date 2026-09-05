package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus string
const (
	AttPresent  AttendanceStatus = "present"
	AttAbsent   AttendanceStatus = "absent"
	AttLate     AttendanceStatus = "late"
	AttHalfDay  AttendanceStatus = "half_day"
	AttHoliday  AttendanceStatus = "holiday"
	AttLeave    AttendanceStatus = "leave"
)

type AttendanceRecord struct {
	BaseEntity
	EmployeeID     uuid.UUID        `json:"employee_id"`
	Date           time.Time        `json:"date"`
	ClockIn        *time.Time       `json:"clock_in"`
	ClockOut       *time.Time       `json:"clock_out"`
	GeoLocationIn  string           `json:"geo_location_in"`
	GeoLocationOut string           `json:"geo_location_out"`
	Status         AttendanceStatus `json:"status"`
	OvertimeHours  float64          `json:"overtime_hours"`
	LateMinutes    int              `json:"late_minutes"`
	Notes          string           `json:"notes"`
}

type Shift struct {
	BaseEntity
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	BranchID  uuid.UUID `json:"branch_id"`
	IsActive  bool      `json:"is_active"`
}

type ShiftAssignment struct {
	BaseEntity
	EmployeeID uuid.UUID `json:"employee_id"`
	ShiftID    uuid.UUID `json:"shift_id"`
	Date       time.Time `json:"date"`
}

type ShiftSwapRequest struct {
	BaseEntity
	RequestorID          uuid.UUID `json:"requestor_id"`
	TargetEmployeeID     uuid.UUID `json:"target_employee_id"`
	OriginalAssignmentID uuid.UUID `json:"original_assignment_id"`
	TargetAssignmentID   uuid.UUID `json:"target_assignment_id"`
	Date                 time.Time `json:"date"`
	Status               string    `json:"status"`
	Notes                string    `json:"notes"`
}

type AttendanceRepository interface {
	ClockIn(ctx context.Context, req *AttendanceRecord) error
	ClockOut(ctx context.Context, id uuid.UUID, outTime time.Time, geo string) error
	GetByDateRange(ctx context.Context, empID uuid.UUID, start, end time.Time) ([]AttendanceRecord, error)
}

type AttendanceUsecase interface {
	ClockIn(ctx context.Context, empID uuid.UUID, geo string) (AttendanceRecord, error)
	ClockOut(ctx context.Context, empID uuid.UUID, geo string) (AttendanceRecord, error)
}
