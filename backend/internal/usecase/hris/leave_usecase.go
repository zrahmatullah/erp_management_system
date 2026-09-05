package hris

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type LeaveUsecaseImpl struct {
	repo domain.LeaveRepository
}

func NewLeaveUsecase(repo domain.LeaveRepository) domain.LeaveUsecase {
	return &LeaveUsecaseImpl{repo: repo}
}

func (u *LeaveUsecaseImpl) SubmitRequest(ctx context.Context, req domain.LeaveRequest) (domain.LeaveRequest, error) {
	req.Status = domain.LeavePending
	err := u.repo.CreateRequest(ctx, &req)
	return req, err
}

func (u *LeaveUsecaseImpl) ApproveRequest(ctx context.Context, reqID, approverID uuid.UUID, status domain.LeaveStatus) error {
	req := &domain.LeaveRequest{
		BaseEntity: domain.BaseEntity{ID: reqID},
		Status:     status,
		ApprovedBy: &approverID,
	}
	return u.repo.UpdateRequest(ctx, req)
}
