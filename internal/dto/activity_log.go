package dto

import (
	"POS-kasir/internal/repository"
	"time"

	"github.com/google/uuid"
)

type GetActivityLogsRequest struct {
	Page       int                      `query:"page" validate:"min=1"`
	Limit      int                      `query:"limit" validate:"min=1,max=100"`
	Search     string                   `query:"search"`
	StartDate  string                   `query:"start_date"`
	EndDate    string                   `query:"end_date"`
	UserID     string                   `query:"user_id"`
	EntityType repository.LogEntityType `query:"entity_type"`
	ActionType repository.LogActionType `query:"action_type"`
}

type ActivityLogResponse struct {
	ID         uuid.UUID                `json:"id"`
	UserID     uuid.UUID                `json:"user_id"`
	UserName   string                   `json:"user_name"`
	ActionType repository.LogActionType `json:"action_type"`
	EntityType repository.LogEntityType `json:"entity_type"`
	EntityID   string                   `json:"entity_id"`
	Details    map[string]interface{}   `json:"details"`
	CreatedAt  time.Time                `json:"created_at"`
}

type ActivityLogListResponse struct {
	Logs       []ActivityLogResponse `json:"logs"`
	TotalItems int64                 `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
}
