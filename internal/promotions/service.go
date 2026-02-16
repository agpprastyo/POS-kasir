package promotions

import (
	"POS-kasir/internal/activitylog"
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/common/store"
	"POS-kasir/internal/promotions/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type IPromotionService interface {
	CreatePromotion(ctx context.Context, req CreatePromotionRequest) (*PromotionResponse, error)
	UpdatePromotion(ctx context.Context, id uuid.UUID, req UpdatePromotionRequest) (*PromotionResponse, error)
	DeletePromotion(ctx context.Context, id uuid.UUID) error
	GetPromotion(ctx context.Context, id uuid.UUID) (*PromotionResponse, error)
	ListPromotions(ctx context.Context, req ListPromotionsRequest) (*PagedPromotionResponse, error)
	RestorePromotion(ctx context.Context, id uuid.UUID) error
}

type PromotionService struct {
	repo            repository.Querier
	store           store.Store
	log             logger.ILogger
	activityService activitylog.IActivityService
}

func NewPromotionService(store store.Store, repo repository.Querier, log logger.ILogger) IPromotionService {
	return &PromotionService{
		repo:  repo,
		store: store,
		log:   log,
	}
}

func (s *PromotionService) CreatePromotion(ctx context.Context, req CreatePromotionRequest) (*PromotionResponse, error) {
	var promoID uuid.UUID

	err := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := repository.New(tx)

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

	actorID := ctx.Value("user_id").(uuid.UUID)

	logDetails := map[string]interface{}{
		"created_promotion_id":   promoID,
		"created_promotion_name": req.Name,
	}

	s.activityService.Log(
		ctx,
		actorID,
		activitylog_repo.LogActionTypeCREATE,
		activitylog_repo.LogEntityTypePROMOTION,
		promoID.String(),
		logDetails,
	)

	return s.GetPromotion(ctx, promoID)
}

func (s *PromotionService) UpdatePromotion(ctx context.Context, id uuid.UUID, req UpdatePromotionRequest) (*PromotionResponse, error) {
	err := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := repository.New(tx)
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

	actorID := ctx.Value("user_id").(uuid.UUID)

	s.activityService.Log(
		ctx,
		actorID,
		activitylog_repo.LogActionTypeUPDATE,
		activitylog_repo.LogEntityTypePROMOTION,
		id.String(),
		map[string]interface{}{
			"updated_promotion_id":   id,
			"updated_promotion_name": req.Name,
		},
	)

	return s.GetPromotion(ctx, id)
}

func (s *PromotionService) DeletePromotion(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeletePromotion(ctx, id)
	if err != nil {
		s.log.Error("Failed to delete promotion", "error", err, "id", id)
		return err
	}
	return nil
}

func (s *PromotionService) GetPromotion(ctx context.Context, id uuid.UUID) (*PromotionResponse, error) {
	promo, err := s.repo.GetPromotionByID(ctx, id)
	if err != nil {
		return nil, common.ErrNotFound
	}

	rules, err := s.repo.GetPromotionRules(ctx, id)
	if err != nil {
		return nil, err
	}

	targets, err := s.repo.GetPromotionTargets(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.mapToDetailResponse(promo, rules, targets), nil
}

func (s *PromotionService) ListPromotions(ctx context.Context, req ListPromotionsRequest) (*PagedPromotionResponse, error) {
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
		totalCount, err = s.repo.CountTrashPromotions(ctx)
		if err != nil {
			return nil, err
		}
		promos, err = s.repo.ListTrashPromotions(ctx, repository.ListTrashPromotionsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	} else {
		totalCount, err = s.repo.CountPromotions(ctx)
		if err != nil {
			return nil, err
		}
		promos, err = s.repo.ListPromotions(ctx, repository.ListPromotionsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	}

	if err != nil {
		return nil, err
	}

	var promoResponses []PromotionResponse
	for _, p := range promos {
		rules, _ := s.repo.GetPromotionRules(ctx, p.ID)
		targets, _ := s.repo.GetPromotionTargets(ctx, p.ID)
		promoResponses = append(promoResponses, *s.mapToDetailResponse(p, rules, targets))
	}

	return &PagedPromotionResponse{
		Promotions: promoResponses,
		Pagination: pagination.BuildPagination(page, int(totalCount), limit),
	}, nil
}

func (s *PromotionService) RestorePromotion(ctx context.Context, id uuid.UUID) error {
	err := s.repo.RestorePromotion(ctx, id)
	if err != nil {
		s.log.Error("Failed to restore promotion", "error", err, "id", id)
		return err
	}

	actorID := ctx.Value("user_id").(uuid.UUID)

	s.activityService.Log(
		ctx,
		actorID,
		activitylog_repo.LogActionTypeDELETE,
		activitylog_repo.LogEntityTypePROMOTION,
		id.String(),
		map[string]interface{}{
			"restored_promotion_id": id,
		},
	)
	return nil
}

func (s *PromotionService) mapToDetailResponse(
	p repository.Promotion,
	rules []repository.PromotionRule,
	targets []repository.PromotionTarget,
) *PromotionResponse {

	ruleResponses := make([]PromotionRuleResponse, len(rules))
	for i, r := range rules {
		var desc string
		if r.Description != nil {
			desc = *r.Description
		}
		ruleResponses[i] = PromotionRuleResponse{
			ID:          r.ID,
			RuleType:    r.RuleType,
			RuleValue:   r.RuleValue,
			Description: desc,
		}
	}

	targetResponses := make([]PromotionTargetResponse, len(targets))
	for i, t := range targets {
		targetResponses[i] = PromotionTargetResponse{
			ID:         t.ID,
			TargetType: t.TargetType,
			TargetID:   t.TargetID,
		}
	}

	var desc string
	if p.Description != nil {
		desc = *p.Description
	}

	return &PromotionResponse{
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
