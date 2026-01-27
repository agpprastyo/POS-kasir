package dto

import (
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/repository"
	"time"

	"github.com/google/uuid"
)

type CreatePromotionRuleRequest struct {
	RuleType    repository.PromotionRuleType `json:"rule_type" validate:"required"`
	RuleValue   string                       `json:"rule_value" validate:"required"`
	Description string                       `json:"description"`
}

type CreatePromotionTargetRequest struct {
	TargetType repository.PromotionTargetType `json:"target_type" validate:"required"`
	TargetID   string                         `json:"target_id" validate:"required"`
}

type CreatePromotionRequest struct {
	Name              string                         `json:"name" validate:"required,min=3"`
	Description       string                         `json:"description"`
	Scope             repository.PromotionScope      `json:"scope" validate:"required"`
	DiscountType      repository.DiscountType        `json:"discount_type" validate:"required"`
	DiscountValue     int64                          `json:"discount_value" validate:"required,gt=0"`
	MaxDiscountAmount *int64                         `json:"max_discount_amount"`
	StartDate         time.Time                      `json:"start_date" validate:"required"`
	EndDate           time.Time                      `json:"end_date" validate:"required,gtfield=StartDate"`
	IsActive          bool                           `json:"is_active"`
	Rules             []CreatePromotionRuleRequest   `json:"rules" validate:"dive"`
	Targets           []CreatePromotionTargetRequest `json:"targets" validate:"dive"`
}

type UpdatePromotionRequest struct {
	Name              string                         `json:"name" validate:"required,min=3"`
	Description       string                         `json:"description"`
	Scope             repository.PromotionScope      `json:"scope" validate:"required"`
	DiscountType      repository.DiscountType        `json:"discount_type" validate:"required"`
	DiscountValue     int64                          `json:"discount_value" validate:"required,gt=0"`
	MaxDiscountAmount *int64                         `json:"max_discount_amount"`
	StartDate         time.Time                      `json:"start_date" validate:"required"`
	EndDate           time.Time                      `json:"end_date" validate:"required,gtfield=StartDate"`
	IsActive          bool                           `json:"is_active"`
	Rules             []CreatePromotionRuleRequest   `json:"rules" validate:"dive"`
	Targets           []CreatePromotionTargetRequest `json:"targets" validate:"dive"`
}

type PromotionRuleResponse struct {
	ID          uuid.UUID                    `json:"id"`
	RuleType    repository.PromotionRuleType `json:"rule_type"`
	RuleValue   string                       `json:"rule_value"`
	Description string                       `json:"description,omitempty"`
}

type PromotionTargetResponse struct {
	ID         uuid.UUID                      `json:"id"`
	TargetType repository.PromotionTargetType `json:"target_type"`
	TargetID   string                         `json:"target_id"`
}

type PromotionResponse struct {
	ID                uuid.UUID                 `json:"id"`
	Name              string                    `json:"name"`
	Description       string                    `json:"description,omitempty"`
	Scope             repository.PromotionScope `json:"scope"`
	DiscountType      repository.DiscountType   `json:"discount_type"`
	DiscountValue     int64                     `json:"discount_value"`
	MaxDiscountAmount *int64                    `json:"max_discount_amount,omitempty"`
	StartDate         time.Time                 `json:"start_date"`
	EndDate           time.Time                 `json:"end_date"`
	IsActive          bool                      `json:"is_active"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	DeletedAt         *time.Time                `json:"deleted_at,omitempty"`
	Rules             []PromotionRuleResponse   `json:"rules"`
	Targets           []PromotionTargetResponse `json:"targets"`
}

type ListPromotionsRequest struct {
	Page  *int `query:"page"`
	Limit *int `query:"limit"`
	Trash bool `query:"trash"`
}

type PagedPromotionResponse struct {
	Promotions []PromotionResponse   `json:"promotions"`
	Pagination pagination.Pagination `json:"pagination"`
}
