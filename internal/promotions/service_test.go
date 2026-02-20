package promotions_test

import (
	"POS-kasir/internal/activitylog/repository"
	common "POS-kasir/internal/common"
	"POS-kasir/internal/promotions"
	promo_repo "POS-kasir/internal/promotions/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

func TestPromotionService_CreatePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockActivityService := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	// Mock DB for s.repo AND transaction generation
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Create a legitimate pgx.Tx object from the mock pool
	// We need to expect the Begin call first
	mockDB.ExpectBegin()
	mockTx, err := mockDB.Begin(context.Background())
	assert.NoError(t, err)

	repo := promo_repo.New(mockDB)
	service := promotions.NewPromotionService(mockStore, repo, mockLogger, mockActivityService)

	ctx := context.WithValue(context.Background(), "user_id", uuid.New())
	userID := ctx.Value("user_id").(uuid.UUID)

	req := promotions.CreatePromotionRequest{
		Name:          "Promo Merdeka",
		Description:   "Diskon Kemerdekaan",
		Scope:         promo_repo.PromotionScopeORDER,
		DiscountType:  promo_repo.DiscountTypePercentage,
		DiscountValue: 10,
		StartDate:     time.Now(),
		EndDate:       time.Now().Add(24 * time.Hour),
		IsActive:      true,
		Rules: []promotions.CreatePromotionRuleRequest{
			{
				RuleType:  promo_repo.PromotionRuleTypeMINIMUMORDERAMOUNT,
				RuleValue: "50000",
			},
		},
		Targets: []promotions.CreatePromotionTargetRequest{}, // No targets for Order scope example
	}

	promoID := uuid.New()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		// Expect ExecTx to call the function with mockTx
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			return fn(mockTx)
		})

		// Expect INSERT Promotion on mockTx (recorded on mockDB)
		mockDB.ExpectQuery("INSERT INTO promotions").
			WithArgs(
				req.Name,
				&req.Description,
				req.Scope,
				req.DiscountType,
				utils.Int64ToNumeric(req.DiscountValue),
				utils.Int64PtrToNumeric(req.MaxDiscountAmount),
				pgtype.Timestamptz{Time: req.StartDate, Valid: true},
				pgtype.Timestamptz{Time: req.EndDate, Valid: true},
				req.IsActive,
			).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "scope", "discount_type", "discount_value", "max_discount_amount", "start_date", "end_date", "is_active", "created_at", "updated_at", "deleted_at"}).
				AddRow(promoID, req.Name, &req.Description, req.Scope, req.DiscountType, utils.Int64ToNumeric(req.DiscountValue), utils.Int64PtrToNumeric(req.MaxDiscountAmount), pgtype.Timestamptz{Time: req.StartDate, Valid: true}, pgtype.Timestamptz{Time: req.EndDate, Valid: true}, req.IsActive, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{}))

		// Expect INSERT Promotion Rule on mockTx (recorded on mockDB)
		mockDB.ExpectQuery("INSERT INTO promotion_rules").
			WithArgs(
				promoID,
				req.Rules[0].RuleType,
				req.Rules[0].RuleValue,
				(*string)(nil),
			).
			WillReturnRows(pgxmock.NewRows([]string{"id", "promotion_id", "rule_type", "rule_value", "description", "created_at", "updated_at"}).
				AddRow(uuid.New(), promoID, req.Rules[0].RuleType, req.Rules[0].RuleValue, nil, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}))

		// Expect Activity Log (on successful commit)
		mockActivityService.EXPECT().Log(ctx, userID, repository.LogActionTypeCREATE, repository.LogEntityTypePROMOTION, promoID.String(), gomock.Any())

		// Expect GetPromotion queries on mockDB
		// 1. GetPromotionByID
		mockDB.ExpectQuery("SELECT id, name, description, scope, discount_type, discount_value, max_discount_amount, start_date, end_date, is_active, created_at, updated_at, deleted_at FROM promotions WHERE id = \\$1 LIMIT 1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "scope", "discount_type", "discount_value", "max_discount_amount", "start_date", "end_date", "is_active", "created_at", "updated_at", "deleted_at"}).
				AddRow(promoID, req.Name, &req.Description, req.Scope, req.DiscountType, utils.Int64ToNumeric(req.DiscountValue), utils.Int64PtrToNumeric(req.MaxDiscountAmount), pgtype.Timestamptz{Time: req.StartDate, Valid: true}, pgtype.Timestamptz{Time: req.EndDate, Valid: true}, req.IsActive, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{}))

		// 2. GetPromotionRules
		mockDB.ExpectQuery("SELECT id, promotion_id, rule_type, rule_value, description, created_at, updated_at FROM promotion_rules WHERE promotion_id = \\$1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "promotion_id", "rule_type", "rule_value", "description", "created_at", "updated_at"}).
				AddRow(uuid.New(), promoID, req.Rules[0].RuleType, req.Rules[0].RuleValue, nil, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}))

		// 3. GetPromotionTargets
		mockDB.ExpectQuery("SELECT id, promotion_id, target_type, target_id, created_at, updated_at FROM promotion_targets WHERE promotion_id = \\$1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "promotion_id", "target_type", "target_id", "created_at", "updated_at"})) // Empty targets

		resp, err := service.CreatePromotion(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Name, resp.Name)
		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("TxFailure", func(t *testing.T) {
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).Return(errors.New("db error"))

		// Logger error expectation
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.CreatePromotion(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestPromotionService_UpdatePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockActivityService := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	mockDB.ExpectBegin()
	mockTx, err := mockDB.Begin(context.Background())
	assert.NoError(t, err)

	repo := promo_repo.New(mockDB)
	service := promotions.NewPromotionService(mockStore, repo, mockLogger, mockActivityService)

	ctx := context.WithValue(context.Background(), "user_id", uuid.New())
	userID := ctx.Value("user_id").(uuid.UUID)
	promoID := uuid.New()

	req := promotions.UpdatePromotionRequest{
		Name:          "Promo Update",
		Description:   "Updated Desc",
		Scope:         promo_repo.PromotionScopeITEM,
		DiscountType:  promo_repo.DiscountTypeFixedAmount,
		DiscountValue: 5000,
		StartDate:     time.Now(),
		EndDate:       time.Now().Add(48 * time.Hour),
		IsActive:      true,
		Rules:         []promotions.CreatePromotionRuleRequest{},
		Targets:       []promotions.CreatePromotionTargetRequest{},
	}
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			return fn(mockTx)
		})

		// Update Promotion
		mockDB.ExpectQuery("UPDATE promotions").
			WithArgs(
				promoID,
				req.Name,
				&req.Description,
				req.Scope,
				req.DiscountType,
				utils.Int64ToNumeric(req.DiscountValue),
				utils.Int64PtrToNumeric(req.MaxDiscountAmount),
				pgtype.Timestamptz{Time: req.StartDate, Valid: true},
				pgtype.Timestamptz{Time: req.EndDate, Valid: true},
				req.IsActive,
			).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "scope", "discount_type", "discount_value", "max_discount_amount", "start_date", "end_date", "is_active", "created_at", "updated_at", "deleted_at"}).
				AddRow(promoID, req.Name, &req.Description, req.Scope, req.DiscountType, utils.Int64ToNumeric(req.DiscountValue), utils.Int64PtrToNumeric(req.MaxDiscountAmount), pgtype.Timestamptz{Time: req.StartDate, Valid: true}, pgtype.Timestamptz{Time: req.EndDate, Valid: true}, req.IsActive, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{}))

		// Delete old rules/targets
		mockDB.ExpectExec("DELETE FROM promotion_rules").WithArgs(promoID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mockDB.ExpectExec("DELETE FROM promotion_targets").WithArgs(promoID).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		// New rules/targets (empty in request)

		mockActivityService.EXPECT().Log(ctx, userID, repository.LogActionTypeUPDATE, repository.LogEntityTypePROMOTION, promoID.String(), gomock.Any())

		// GetPromotion queries on mockDB
		mockDB.ExpectQuery("SELECT id, name, description, scope, discount_type, discount_value, max_discount_amount, start_date, end_date, is_active, created_at, updated_at, deleted_at FROM promotions WHERE id = \\$1 LIMIT 1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "scope", "discount_type", "discount_value", "max_discount_amount", "start_date", "end_date", "is_active", "created_at", "updated_at", "deleted_at"}).
				AddRow(promoID, req.Name, &req.Description, req.Scope, req.DiscountType, utils.Int64ToNumeric(req.DiscountValue), utils.Int64PtrToNumeric(req.MaxDiscountAmount), pgtype.Timestamptz{Time: req.StartDate, Valid: true}, pgtype.Timestamptz{Time: req.EndDate, Valid: true}, req.IsActive, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{}))

		mockDB.ExpectQuery("SELECT id, promotion_id, rule_type, rule_value, description, created_at, updated_at FROM promotion_rules").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{}))

		mockDB.ExpectQuery("SELECT id, promotion_id, target_type, target_id, created_at, updated_at FROM promotion_targets").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{}))

		resp, err := service.UpdatePromotion(ctx, promoID, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("TxFailure", func(t *testing.T) {
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).Return(errors.New("tx error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.UpdatePromotion(ctx, promoID, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockStore.EXPECT().ExecTx(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			mockDB.ExpectBegin()
			tx, _ := mockDB.Begin(ctx)
			defer tx.Rollback(ctx)

			// Expect Update Query to fail
			mockDB.ExpectQuery("UPDATE promotions").WithArgs(gomock.Any()).WillReturnError(errors.New("update error"))

			return fn(tx)
		})
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.UpdatePromotion(ctx, promoID, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestPromotionService_ListPromotions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockActivityService := mocks.NewMockIActivityService(ctrl)
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := promo_repo.New(mockDB)
	service := promotions.NewPromotionService(mockStore, repo, mockLogger, mockActivityService)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		page := 1
		limit := 10

		mockDB.ExpectQuery("SELECT COUNT\\(\\*\\) FROM promotions WHERE deleted_at IS NULL").
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(1)))

		mockDB.ExpectQuery("SELECT id, name, description, scope, discount_type, discount_value, max_discount_amount, start_date, end_date, is_active, created_at, updated_at, deleted_at FROM promotions WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \\$1 OFFSET \\$2").
			WithArgs(int32(limit), int32(0)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "scope", "discount_type", "discount_value", "max_discount_amount", "start_date", "end_date", "is_active", "created_at", "updated_at", "deleted_at"}).
				AddRow(uuid.New(), "Promo 1", utils.StringPtr("Desc"), promo_repo.PromotionScopeORDER, promo_repo.DiscountTypePercentage, utils.Int64ToNumeric(10), utils.Int64ToNumeric(5000), pgtype.Timestamptz{Time: time.Now(), Valid: true}, pgtype.Timestamptz{Time: time.Now(), Valid: true}, true, pgtype.Timestamptz{Time: time.Now(), Valid: true}, pgtype.Timestamptz{Time: time.Now(), Valid: true}, pgtype.Timestamptz{}))

		// Get rules and targets for the promotion
		mockDB.ExpectQuery("SELECT id, promotion_id, rule_type, rule_value, description, created_at, updated_at FROM promotion_rules").
			WillReturnRows(pgxmock.NewRows([]string{}))
		mockDB.ExpectQuery("SELECT id, promotion_id, target_type, target_id, created_at, updated_at FROM promotion_targets").
			WithArgs(mock.Anything). // Hard to predict ID here without capturing it, but regex match on query is enough usually. Or match arg type.
			WillReturnRows(pgxmock.NewRows([]string{}))

		resp, err := service.ListPromotions(ctx, promotions.ListPromotionsRequest{Page: &page, Limit: &limit})
		assert.NoError(t, err)
		assert.Len(t, resp.Promotions, 1)
	})
}

func TestPromotionService_DeletePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockActivityService := mocks.NewMockIActivityService(ctrl)
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := promo_repo.New(mockDB)
	service := promotions.NewPromotionService(mockStore, repo, mockLogger, mockActivityService)
	ctx := context.Background()
	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectExec("UPDATE promotions SET deleted_at = NOW\\(\\) WHERE id = \\$1").
			WithArgs(promoID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := service.DeletePromotion(ctx, promoID)
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		mockDB.ExpectExec("UPDATE promotions SET deleted_at = NOW\\(\\) WHERE id = \\$1").
			WithArgs(promoID).
			WillReturnError(errors.New("delete error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		err := service.DeletePromotion(ctx, promoID)
		assert.Error(t, err)
	})
}

func TestPromotionService_GetPromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockActivityService := mocks.NewMockIActivityService(ctrl)
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := promo_repo.New(mockDB)
	service := promotions.NewPromotionService(mockStore, repo, mockLogger, mockActivityService)
	ctx := context.Background()
	promoID := uuid.New()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT id, name, description, scope, discount_type, discount_value, max_discount_amount, start_date, end_date, is_active, created_at, updated_at, deleted_at FROM promotions WHERE id = \\$1 LIMIT 1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "scope", "discount_type", "discount_value", "max_discount_amount", "start_date", "end_date", "is_active", "created_at", "updated_at", "deleted_at"}).
				AddRow(promoID, "Promo Get", utils.StringPtr("Desc"), promo_repo.PromotionScopeORDER, promo_repo.DiscountTypePercentage, utils.Int64ToNumeric(10), nil, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}, true, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{}))

		mockDB.ExpectQuery("SELECT id, promotion_id, rule_type, rule_value, description, created_at, updated_at FROM promotion_rules WHERE promotion_id = \\$1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "promotion_id", "rule_type", "rule_value", "description", "created_at", "updated_at"}))

		mockDB.ExpectQuery("SELECT id, promotion_id, target_type, target_id, created_at, updated_at FROM promotion_targets WHERE promotion_id = \\$1").
			WithArgs(promoID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "promotion_id", "target_type", "target_id", "created_at", "updated_at"}))

		resp, err := service.GetPromotion(ctx, promoID)
		assert.NoError(t, err)
		assert.Equal(t, "Promo Get", resp.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDB.ExpectQuery("SELECT id, name, description, scope, discount_type, discount_value, max_discount_amount, start_date, end_date, is_active, created_at, updated_at, deleted_at FROM promotions WHERE id = \\$1 LIMIT 1").
			WithArgs(promoID).
			WillReturnError(pgx.ErrNoRows)

		resp, err := service.GetPromotion(ctx, promoID)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, errors.Is(err, common.ErrNotFound))
	})
}

func TestPromotionService_RestorePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockActivityService := mocks.NewMockIActivityService(ctrl)
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := promo_repo.New(mockDB)
	service := promotions.NewPromotionService(mockStore, repo, mockLogger, mockActivityService)
	ctx := context.WithValue(context.Background(), "user_id", uuid.New())
	userID := ctx.Value("user_id").(uuid.UUID)
	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockDB.ExpectExec("UPDATE promotions SET deleted_at = NULL WHERE id = \\$1").
			WithArgs(promoID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		mockActivityService.EXPECT().Log(ctx, userID, repository.LogActionTypeRESTORE, repository.LogEntityTypePROMOTION, promoID.String(), gomock.Any())

		err := service.RestorePromotion(ctx, promoID)
		assert.NoError(t, err)
	})
}
