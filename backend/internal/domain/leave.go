package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LeaveStatus string
const (
	LeavePending         LeaveStatus = "pending"
	LeaveManagerApproved LeaveStatus = "manager_approved"
	LeaveHRApproved      LeaveStatus = "hr_approved"
	LeaveApproved        LeaveStatus = "approved"
	LeaveRejected        LeaveStatus = "rejected"
	LeaveCancelled       LeaveStatus = "cancelled"
)

type LeaveType struct {
	BaseEntity
	Name        string `json:"name"`
	DefaultDays int    `json:"default_days"`
	IsPaid      bool   `json:"is_paid"`
	Description string `json:"description"`
}

type LeaveBalance struct {
	BaseEntity
	EmployeeID  uuid.UUID `json:"employee_id"`
	LeaveTypeID uuid.UUID `json:"leave_type_id"`
	Year        int       `json:"year"`
	Allocated   int       `json:"allocated"`
	Used        int       `json:"used"`
	Remaining   int       `json:"remaining"`
}

type LeaveRequest struct {
	BaseEntity
	EmployeeID    uuid.UUID   `json:"employee_id"`
	LeaveTypeID   uuid.UUID   `json:"leave_type_id"`
	StartDate     time.Time   `json:"start_date"`
	EndDate       time.Time   `json:"end_date"`
	TotalDays     int         `json:"total_days"`
	Reason        string      `json:"reason"`
	Status        LeaveStatus `json:"status"`
	ApprovedBy    *uuid.UUID  `json:"approved_by"`
	ApproverNotes string      `json:"approver_notes"`
}

type LeaveRepository interface {
	CreateRequest(ctx context.Context, req *LeaveRequest) error
	UpdateRequest(ctx context.Context, req *LeaveRequest) error
	GetBalance(ctx context.Context, empID uuid.UUID, year int) ([]LeaveBalance, error)
}

type LeaveUsecase interface {
	SubmitRequest(ctx context.Context, req LeaveRequest) (LeaveRequest, error)
	ApproveRequest(ctx context.Context, reqID, approverID uuid.UUID, status LeaveStatus) error
}
