package activitylog

import (
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"time"

	"github.com/google/uuid"
)

type GetActivityLogsRequest struct {
	Page       *int                            `query:"page" validate:"omitempty,min=1"`
	Limit      *int                            `query:"limit" validate:"omitempty,min=1,max=100"`
	Search     *string                         `query:"search" validate:"omitempty"`
	StartDate  *string                         `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    *string                         `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	UserID     *string                         `query:"user_id" validate:"omitempty,uuid"`
	EntityType *activitylog_repo.LogEntityType `query:"entity_type" validate:"omitempty,oneof=PRODUCT CATEGORY PROMOTION ORDER USER"`
	ActionType *activitylog_repo.LogActionType `query:"action_type" validate:"omitempty,oneof=CREATE UPDATE DELETE CANCEL APPLY_PROMOTION PROCESS_PAYMENT REGISTER UPDATE_PASSWORD UPDATE_AVATAR LOGIN_SUCCESS LOGIN_FAILED"`
}

type ActivityLogResponse struct {
	ID         uuid.UUID                      `json:"id"`
	UserID     uuid.UUID                      `json:"user_id"`
	UserName   string                         `json:"user_name"`
	ActionType activitylog_repo.LogActionType `json:"action_type"`
	EntityType activitylog_repo.LogEntityType `json:"entity_type"`
	EntityID   string                         `json:"entity_id"`
	Details    map[string]interface{}         `json:"details"`
	CreatedAt  time.Time                      `json:"created_at"`
}

type ActivityLogListResponse struct {
	Logs       []ActivityLogResponse `json:"logs"`
	TotalItems int64                 `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
}
