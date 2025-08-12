package dto

import "time"

type CancellationReasonResponse struct {
	ID          int32     `json:"id"`
	Reason      string    `json:"reason"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}
