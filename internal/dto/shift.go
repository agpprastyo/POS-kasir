package dto

import (
	"POS-kasir/internal/repository"
	"time"

	"github.com/google/uuid"
)

type StartShiftRequest struct {
	StartCash int64  `json:"start_cash" validate:"min=0"`
	Password  string `json:"password" validate:"required"`
}

type EndShiftRequest struct {
	ActualCashEnd int64  `json:"actual_cash_end" validate:"min=0"`
	Password      string `json:"password" validate:"required"`
}

type ShiftResponse struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         *time.Time             `json:"end_time,omitempty"`
	StartCash       int64                  `json:"start_cash"`
	ExpectedCashEnd *int64                 `json:"expected_cash_end,omitempty"`
	ActualCashEnd   *int64                 `json:"actual_cash_end,omitempty"`
	Difference      *int64                 `json:"difference,omitempty"` // Actual - Expected
	Status          repository.ShiftStatus `json:"status"`
}

type CashTransactionRequest struct {
	Amount      int64                          `json:"amount" validate:"required,min=1"`
	Type        repository.CashTransactionType `json:"type" validate:"required,oneof=cash_in cash_out"`
	Category    string                         `json:"category" validate:"required"`
	Description string                         `json:"description"`
}

type CashTransactionResponse struct {
	ID          uuid.UUID                      `json:"id"`
	ShiftID     uuid.UUID                      `json:"shift_id"`
	UserID      uuid.UUID                      `json:"user_id"`
	Amount      int64                          `json:"amount"`
	Type        repository.CashTransactionType `json:"type"`
	Category    string                         `json:"category"`
	Description *string                        `json:"description,omitempty"`
	CreatedAt   time.Time                      `json:"created_at"`
}
