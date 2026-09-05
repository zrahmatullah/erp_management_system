package hris

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type AttendanceUsecaseImpl struct {
	repo domain.AttendanceRepository
}

func NewAttendanceUsecase(repo domain.AttendanceRepository) domain.AttendanceUsecase {
	return &AttendanceUsecaseImpl{repo: repo}
}

func (u *AttendanceUsecaseImpl) ClockIn(ctx context.Context, empID uuid.UUID, geo string) (domain.AttendanceRecord, error) {
	now := time.Now()
	
	// Check if already clocked in today (omitted for brevity)

	record := domain.AttendanceRecord{
		EmployeeID:    empID,
		Date:          time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		ClockIn:       &now,
		GeoLocationIn: geo,
		Status:        domain.AttPresent,
	}

	err := u.repo.ClockIn(ctx, &record)
	return record, err
}

func (u *AttendanceUsecaseImpl) ClockOut(ctx context.Context, attID uuid.UUID, geo string) (domain.AttendanceRecord, error) {
	if attID == uuid.Nil {
		return domain.AttendanceRecord{}, errors.New("invalid attendance record id")
	}

	now := time.Now()
	err := u.repo.ClockOut(ctx, attID, now, geo)
	if err != nil {
		return domain.AttendanceRecord{}, err
	}

	return domain.AttendanceRecord{}, nil // fetch and return real record in practice
}
