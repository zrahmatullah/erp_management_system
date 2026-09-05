package notification

import (
	"context"

	"github.com/google/uuid"

	"cafe-erp-system/backend/internal/domain"
)

type NotificationUsecaseImpl struct {
	repo domain.NotificationRepository
}

func NewNotificationUsecase(repo domain.NotificationRepository) *NotificationUsecaseImpl {
	return &NotificationUsecaseImpl{repo: repo}
}

func (u *NotificationUsecaseImpl) CreateNotification(ctx context.Context, notif domain.Notification) error {
	return u.repo.Create(ctx, &notif)
}

func (u *NotificationUsecaseImpl) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return u.repo.MarkAsRead(ctx, id)
}
