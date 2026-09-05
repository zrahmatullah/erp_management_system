package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"cafe-erp-system/backend/internal/domain"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) domain.NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, title, content, type, reference_type, reference_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING created_at`

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		n.ID, n.UserID, n.Title, n.Content, n.Type, n.ReferenceType, n.ReferenceID, n.IsRead,
	).Scan(&n.CreatedAt)
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

type AuditLogRepository struct {
	db *pgxpool.Pool
}

func NewAuditLogRepository(db *pgxpool.Pool) domain.AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Log(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, table_name, record_id, old_data, new_data, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())`

	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		log.ID, log.UserID, log.Action, log.TableName, log.RecordID, log.OldData, log.NewData, log.IPAddress, log.UserAgent,
	)
	return err
}
