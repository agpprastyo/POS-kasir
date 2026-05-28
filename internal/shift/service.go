package shift

import (
	"POS-kasir/internal/common"
	repository "POS-kasir/internal/shift/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/metrics"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	StartShift(ctx context.Context, userID uuid.UUID, req StartShiftRequest) (*ShiftResponse, error)
	EndShift(ctx context.Context, userID uuid.UUID, req EndShiftRequest) (*ShiftResponse, error)
	GetOpenShift(ctx context.Context, userID uuid.UUID) (*ShiftResponse, error)
	CreateCashTransaction(ctx context.Context, userID uuid.UUID, req CashTransactionRequest) (*CashTransactionResponse, error)
	AutoCloseShifts(ctx context.Context) error
}

type service struct {
	repo  repository.Querier
	log   logger.ILogger
	cache *Cache
}

func NewService(repo repository.Querier, log logger.ILogger, cache *Cache) Service {
	return &service{
		repo:  repo,
		log:   log,
		cache: cache,
	}
}

func (s *service) StartShift(ctx context.Context, userID uuid.UUID, req StartShiftRequest) (*ShiftResponse, error) {
	// Verify user password
	passwordHash, err := s.repo.GetUserPasswordHash(ctx, userID)
	if err != nil {
		s.log.Errorf("StartShift | User not found or hash missing: %v", userID)
		return nil, common.ErrNotFound
	}

	if !utils.CheckPassword(passwordHash, req.Password) {
		s.log.Errorf("StartShift | Invalid password for user: %v", userID)
		return nil, common.ErrInvalidCredentials
	}

	// Check if user already has an open shift
	_, err = s.repo.GetOpenShiftByUserID(ctx, userID)
	if err == nil {
		s.log.Warnf("StartShift | User already has an open shift: userID=%v", userID)
		return nil, errors.New("user already has an open shift")
	}

	shift, err := s.repo.CreateShift(ctx, repository.CreateShiftParams{
		UserID:    userID,
		StartCash: req.StartCash,
	})
	if err != nil {
		s.log.Errorf("StartShift | Failed to create shift: %v, userID=%v", err, userID)
		return nil, err
	}

	// Update cache to reflect an open shift
	if s.cache != nil {
		s.cache.SetOpen(userID, true)
	}

	// Record shift opened metric
	metrics.ActiveShifts.Inc()

	return s.mapShiftToResponse(shift), nil
}

func (s *service) EndShift(ctx context.Context, userID uuid.UUID, req EndShiftRequest) (*ShiftResponse, error) {
	// Verify user password
	passwordHash, err := s.repo.GetUserPasswordHash(ctx, userID)
	if err != nil {
		s.log.Errorf("EndShift | User not found or hash missing: %v", userID)
		return nil, common.ErrNotFound
	}

	if !utils.CheckPassword(passwordHash, req.Password) {
		s.log.Errorf("EndShift | Invalid password for user: %v", userID)
		return nil, common.ErrInvalidCredentials
	}

	// Get current open shift
	shift, err := s.repo.GetOpenShiftByUserID(ctx, userID)
	if err != nil {
		s.log.Warnf("EndShift | No open shift found: userID=%v", userID)
		return nil, errors.New("no open shift found")
	}

	// Calculate expected cash end
	// Start Cash + Cash In - Cash Out + Sales (Cash Payment method)
	// For MVP, we will simplify: Start Cash + Cash In - Cash Out.
	// Ideally we need to query orders with payment_method_id corresponding to CASH for this shift.
	// Since we don't have direct link between orders and shifts yet (implied by time and user),
	// we would need to sum order totals where user_id = shift.user_id AND created_at >= shift.start_time AND payment_method is CASH.

	// For now, let's just sum cash transactions. Note: This is partial implementation.
	// TODO: Include Sales in Expected Cash Calculation.
	cashIn, err := s.repo.GetCashTotalByShiftIDAndType(ctx, repository.GetCashTotalByShiftIDAndTypeParams{
		ShiftID: shift.ID,
		Type:    repository.CashTransactionTypeCashIn,
	})
	if err != nil {
		s.log.Errorf("EndShift | Failed to get cash in total: %v", err)
		return nil, err
	}

	cashOut, err := s.repo.GetCashTotalByShiftIDAndType(ctx, repository.GetCashTotalByShiftIDAndTypeParams{
		ShiftID: shift.ID,
		Type:    repository.CashTransactionTypeCashOut,
	})
	if err != nil {
		s.log.Errorf("EndShift | Failed to get cash out total: %v", err)
		return nil, err
	}

	cashSales, err := s.repo.GetCashSalesDuringShift(ctx, repository.GetCashSalesDuringShiftParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		CreatedAt: shift.StartTime,
	})
	if err != nil {
		s.log.Errorf("EndShift | Failed to get cash sales: %v", err)
		cashSales = 0 // don't fail the whole operation
	}

	expectedCashEnd := shift.StartCash + cashIn - cashOut + cashSales

	updatedShift, err := s.repo.EndShift(ctx, repository.EndShiftParams{
		ID:              shift.ID,
		ExpectedCashEnd: &expectedCashEnd,
		ActualCashEnd:   &req.ActualCashEnd,
	})
	if err != nil {
		s.log.Errorf("EndShift | Failed to update shift: %v", err)
		return nil, err
	}

	res := s.mapShiftToResponse(updatedShift)

	diff := req.ActualCashEnd - expectedCashEnd
	res.Difference = &diff

	// Update cache (Clear)
	if s.cache != nil {
		s.cache.Clear(userID)
	}

	// Record shift closed metric
	metrics.ActiveShifts.Dec()

	return res, nil
}

func (s *service) GetOpenShift(ctx context.Context, userID uuid.UUID) (*ShiftResponse, error) {
	// 1. Check cache first
	if s.cache != nil {
		isOpen, found := s.cache.GetOpen(userID)
		if found && !isOpen {
			return nil, common.ErrNotFound
		}
	}

	// 2. Fetch from repository
	shift, err := s.repo.GetOpenShiftByUserID(ctx, userID)
	if err != nil {
		// Update cache if cache is available
		if s.cache != nil {
			s.cache.SetOpen(userID, false)
		}
		return nil, common.ErrNotFound
	}

	// 3. Update cache to "open"
	if s.cache != nil {
		s.cache.SetOpen(userID, true)
	}

	return s.mapShiftToResponse(shift), nil
}

func (s *service) CreateCashTransaction(ctx context.Context, userID uuid.UUID, req CashTransactionRequest) (*CashTransactionResponse, error) {
	// Get current open shift
	shift, err := s.repo.GetOpenShiftByUserID(ctx, userID)
	if err != nil {
		s.log.Warnf("CreateCashTransaction | No open shift found: userID=%v", userID)
		return nil, errors.New("no open shift found")
	}

	tx, err := s.repo.CreateCashTransaction(ctx, repository.CreateCashTransactionParams{
		ShiftID:     shift.ID,
		UserID:      userID,
		Amount:      req.Amount,
		Type:        req.Type,
		Category:    req.Category,
		Description: &req.Description,
	})
	if err != nil {
		s.log.Errorf("CreateCashTransaction | Failed to create transaction: %v", err)
		return nil, err
	}

	return &CashTransactionResponse{
		ID:          tx.ID,
		ShiftID:     tx.ShiftID,
		UserID:      tx.UserID,
		Amount:      tx.Amount,
		Type:        tx.Type,
		Category:    tx.Category,
		Description: tx.Description,
		CreatedAt:   tx.CreatedAt.Time,
	}, nil
}

func (s *service) AutoCloseShifts(ctx context.Context) error {
	shifts, err := s.repo.GetOpenShifts(ctx)
	if err != nil {
		s.log.Errorf("AutoCloseShifts | Failed to get open shifts: %v", err)
		return err
	}

	for _, shift := range shifts {
		s.log.Infof("AutoCloseShifts | Closing shift %v for user %v", shift.ID, shift.UserID)

		// Calculate expected cash
		cashIn, _ := s.repo.GetCashTotalByShiftIDAndType(ctx, repository.GetCashTotalByShiftIDAndTypeParams{
			ShiftID: shift.ID,
			Type:    repository.CashTransactionTypeCashIn,
		})
		cashOut, _ := s.repo.GetCashTotalByShiftIDAndType(ctx, repository.GetCashTotalByShiftIDAndTypeParams{
			ShiftID: shift.ID,
			Type:    repository.CashTransactionTypeCashOut,
		})

		cashSales, _ := s.repo.GetCashSalesDuringShift(ctx, repository.GetCashSalesDuringShiftParams{
			UserID:    pgtype.UUID{Bytes: shift.UserID, Valid: true},
			CreatedAt: shift.StartTime,
		})

		expectedCashEnd := shift.StartCash + cashIn - cashOut + cashSales

		// For auto-close, we assume Actual = Expected to avoid difference.
		// Or we can just leave actual as nil? The schema allows nil.
		// But EndShift in repo sets actual_cash_end.
		_, err := s.repo.EndShift(ctx, repository.EndShiftParams{
			ID:              shift.ID,
			ExpectedCashEnd: &expectedCashEnd,
			ActualCashEnd:   &expectedCashEnd,
		})
		if err != nil {
			s.log.Errorf("AutoCloseShifts | Failed to close shift %v: %v", shift.ID, err)
			continue
		}

		// Update cache
		if s.cache != nil {
			s.cache.Clear(shift.UserID)
		}

		// Record shift auto-closed metric
		metrics.ActiveShifts.Dec()
	}

	return nil
}

func (s *service) mapShiftToResponse(shift repository.Shift) *ShiftResponse {
	var endTime *time.Time
	if shift.EndTime.Valid {
		t := shift.EndTime.Time
		endTime = &t
	}

	var expected *int64
	if shift.ExpectedCashEnd != nil {
		expected = shift.ExpectedCashEnd
	}

	var actual *int64
	if shift.ActualCashEnd != nil {
		actual = shift.ActualCashEnd
	}

	return &ShiftResponse{
		ID:              shift.ID,
		UserID:          shift.UserID,
		StartTime:       shift.StartTime.Time,
		EndTime:         endTime,
		StartCash:       shift.StartCash,
		ExpectedCashEnd: expected,
		ActualCashEnd:   actual,
		Status:          shift.Status,
	}
}
