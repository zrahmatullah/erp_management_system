package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Type          string    `json:"type"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}

type NotificationTemplate struct {
	BaseEntity
	Name      string `json:"name"`
	Channel   string `json:"channel"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Variables string `json:"variables"` // JSON list of variables
}

type AuditLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Action    string    `json:"action"`
	TableName string    `json:"table_name"`
	RecordID  string    `json:"record_id"`
	OldData   string    `json:"old_data"` // JSON
	NewData   string    `json:"new_data"` // JSON
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type SystemSetting struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"` // JSON
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NotificationRepository interface {
	Create(ctx context.Context, notif *Notification) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}

type AuditLogRepository interface {
	Log(ctx context.Context, log *AuditLog) error
}
