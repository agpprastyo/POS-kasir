package promotions

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type IPromotionService interface {
	CreatePromotion(ctx context.Context, req dto.CreatePromotionRequest) (*dto.PromotionResponse, error)
	UpdatePromotion(ctx context.Context, id uuid.UUID, req dto.UpdatePromotionRequest) (*dto.PromotionResponse, error)
	DeletePromotion(ctx context.Context, id uuid.UUID) error
	GetPromotion(ctx context.Context, id uuid.UUID) (*dto.PromotionResponse, error)
	ListPromotions(ctx context.Context, req dto.ListPromotionsRequest) (*dto.PagedPromotionResponse, error)
	RestorePromotion(ctx context.Context, id uuid.UUID) error
}

type PromotionService struct {
	store repository.Store
	log   logger.ILogger
}

func NewPromotionService(store repository.Store, log logger.ILogger) IPromotionService {
	return &PromotionService{
		store: store,
		log:   log,
	}
}

func (s *PromotionService) CreatePromotion(ctx context.Context, req dto.CreatePromotionRequest) (*dto.PromotionResponse, error) {
	var promoID uuid.UUID

	err := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		// 1. Create Promotion
		var description *string
		if req.Description != "" {
			description = &req.Description
		}

		promo, err := qtx.CreatePromotion(ctx, repository.CreatePromotionParams{
			Name:              req.Name,
			Description:       description,
			Scope:             req.Scope,
			DiscountType:      req.DiscountType,
			DiscountValue:     utils.Int64ToNumeric(req.DiscountValue),
			MaxDiscountAmount: utils.Int64PtrToNumeric(req.MaxDiscountAmount),
			StartDate:         pgtype.Timestamptz{Time: req.StartDate, Valid: true},
			EndDate:           pgtype.Timestamptz{Time: req.EndDate, Valid: true},
			IsActive:          req.IsActive,
		})
		if err != nil {
			return err
		}
		promoID = promo.ID

		// 2. Insert Rules
		for _, r := range req.Rules {
			var ruleDesc *string
			if r.Description != "" {
				ruleDesc = &r.Description
			}
			_, err := qtx.CreatePromotionRule(ctx, repository.CreatePromotionRuleParams{
				PromotionID: promo.ID,
				RuleType:    r.RuleType,
				RuleValue:   r.RuleValue,
				Description: ruleDesc,
			})
			if err != nil {
				return err
			}
		}

		// 3. Insert Targets
		for _, t := range req.Targets {
			_, err := qtx.CreatePromotionTarget(ctx, repository.CreatePromotionTargetParams{
				PromotionID: promo.ID,
				TargetType:  t.TargetType,
				TargetID:    t.TargetID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.log.Error("Failed to create promotion", "error", err)
		return nil, err
	}

	return s.GetPromotion(ctx, promoID)
}

func (s *PromotionService) UpdatePromotion(ctx context.Context, id uuid.UUID, req dto.UpdatePromotionRequest) (*dto.PromotionResponse, error) {
	err := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		// 1. Update Promotion
		var description *string
		if req.Description != "" {
			description = &req.Description
		}

		_, err := qtx.UpdatePromotion(ctx, repository.UpdatePromotionParams{
			ID:                id,
			Name:              req.Name,
			Description:       description,
			Scope:             req.Scope,
			DiscountType:      req.DiscountType,
			DiscountValue:     utils.Int64ToNumeric(req.DiscountValue),
			MaxDiscountAmount: utils.Int64PtrToNumeric(req.MaxDiscountAmount),
			StartDate:         pgtype.Timestamptz{Time: req.StartDate, Valid: true},
			EndDate:           pgtype.Timestamptz{Time: req.EndDate, Valid: true},
			IsActive:          req.IsActive,
		})
		if err != nil {
			return err
		}

		// 2. Replace Rules (Delete all then insert)
		if err := qtx.DeletePromotionRulesByPromotionID(ctx, id); err != nil {
			return err
		}
		for _, r := range req.Rules {
			var ruleDesc *string
			if r.Description != "" {
				ruleDesc = &r.Description
			}
			_, err := qtx.CreatePromotionRule(ctx, repository.CreatePromotionRuleParams{
				PromotionID: id,
				RuleType:    r.RuleType,
				RuleValue:   r.RuleValue,
				Description: ruleDesc,
			})
			if err != nil {
				return err
			}
		}

		// 3. Replace Targets (Delete all then insert)
		if err := qtx.DeletePromotionTargetsByPromotionID(ctx, id); err != nil {
			return err
		}
		for _, t := range req.Targets {
			_, err := qtx.CreatePromotionTarget(ctx, repository.CreatePromotionTargetParams{
				PromotionID: id,
				TargetType:  t.TargetType,
				TargetID:    t.TargetID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.log.Error("Failed to update promotion", "error", err, "id", id)
		return nil, err
	}

	return s.GetPromotion(ctx, id)
}

func (s *PromotionService) DeletePromotion(ctx context.Context, id uuid.UUID) error {
	// Soft delete by setting is_active = false
	err := s.store.DeletePromotion(ctx, id)
	if err != nil {
		s.log.Error("Failed to delete promotion", "error", err, "id", id)
		return err
	}
	return nil
}

func (s *PromotionService) GetPromotion(ctx context.Context, id uuid.UUID) (*dto.PromotionResponse, error) {
	promo, err := s.store.GetPromotionByID(ctx, id)
	if err != nil {
		return nil, common.ErrNotFound
	}

	rules, err := s.store.GetPromotionRules(ctx, id)
	if err != nil {
		return nil, err
	}

	targets, err := s.store.GetPromotionTargets(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.mapToDetailResponse(promo, rules, targets), nil
}

func (s *PromotionService) ListPromotions(ctx context.Context, req dto.ListPromotionsRequest) (*dto.PagedPromotionResponse, error) {
	page := 1
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	limit := 10
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	var totalCount int64
	var promos []repository.Promotion
	var err error

	if req.Trash {
		totalCount, err = s.store.CountTrashPromotions(ctx)
		if err != nil {
			return nil, err
		}
		promos, err = s.store.ListTrashPromotions(ctx, repository.ListTrashPromotionsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	} else {
		totalCount, err = s.store.CountPromotions(ctx)
		if err != nil {
			return nil, err
		}
		promos, err = s.store.ListPromotions(ctx, repository.ListPromotionsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	}

	if err != nil {
		return nil, err
	}

	var promoResponses []dto.PromotionResponse
	for _, p := range promos {
		rules, _ := s.store.GetPromotionRules(ctx, p.ID)
		targets, _ := s.store.GetPromotionTargets(ctx, p.ID)
		promoResponses = append(promoResponses, *s.mapToDetailResponse(p, rules, targets))
	}

	return &dto.PagedPromotionResponse{
		Promotions: promoResponses,
		Pagination: pagination.BuildPagination(page, int(totalCount), limit),
	}, nil
}

func (s *PromotionService) RestorePromotion(ctx context.Context, id uuid.UUID) error {
	err := s.store.RestorePromotion(ctx, id)
	if err != nil {
		s.log.Error("Failed to restore promotion", "error", err, "id", id)
		return err
	}
	return nil
}

func (s *PromotionService) mapToDetailResponse(
	p repository.Promotion,
	rules []repository.PromotionRule,
	targets []repository.PromotionTarget,
) *dto.PromotionResponse {

	ruleResponses := make([]dto.PromotionRuleResponse, len(rules))
	for i, r := range rules {
		var desc string
		if r.Description != nil {
			desc = *r.Description
		}
		ruleResponses[i] = dto.PromotionRuleResponse{
			ID:          r.ID,
			RuleType:    r.RuleType,
			RuleValue:   r.RuleValue,
			Description: desc,
		}
	}

	targetResponses := make([]dto.PromotionTargetResponse, len(targets))
	for i, t := range targets {
		targetResponses[i] = dto.PromotionTargetResponse{
			ID:         t.ID,
			TargetType: t.TargetType,
			TargetID:   t.TargetID,
		}
	}

	var desc string
	if p.Description != nil {
		desc = *p.Description
	}

	return &dto.PromotionResponse{
		ID:                p.ID,
		Name:              p.Name,
		Description:       desc,
		Scope:             p.Scope,
		DiscountType:      p.DiscountType,
		DiscountValue:     utils.NumericToInt64(p.DiscountValue),
		MaxDiscountAmount: utils.NumericToInt64Ptr(p.MaxDiscountAmount),
		StartDate:         p.StartDate.Time,
		EndDate:           p.EndDate.Time,
		IsActive:          p.IsActive,
		CreatedAt:         p.CreatedAt.Time,
		UpdatedAt:         p.UpdatedAt.Time,
		Rules:             ruleResponses,
		Targets:           targetResponses,
	}
}
