package domain

import (
	"context"

	"github.com/google/uuid"
)

type PerformanceReview struct {
	BaseEntity
	EmployeeID uuid.UUID `json:"employee_id"`
	ReviewerID uuid.UUID `json:"reviewer_id"`
	Period     string    `json:"period"`
	TotalScore float64   `json:"total_score"`
	Status     string    `json:"status"`
	Feedback   string    `json:"feedback"`
}

type KPIDefinition struct {
	BaseEntity
	DepartmentID uuid.UUID `json:"department_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Weight       float64   `json:"weight"`
	TargetValue  float64   `json:"target_value"`
}

type KPIScore struct {
	BaseEntity
	ReviewID    uuid.UUID `json:"review_id"`
	KPIID       uuid.UUID `json:"kpi_id"`
	Score       float64   `json:"score"`
	ActualValue float64   `json:"actual_value"`
	Notes       string    `json:"notes"`
}

type PerformanceRepository interface {
	CreateReview(ctx context.Context, review *PerformanceReview) error
	GetReview(ctx context.Context, id uuid.UUID) (*PerformanceReview, error)
}
