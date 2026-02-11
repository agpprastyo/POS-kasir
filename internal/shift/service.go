package shift

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	StartShift(ctx context.Context, userID uuid.UUID, req dto.StartShiftRequest) (*dto.ShiftResponse, error)
	EndShift(ctx context.Context, userID uuid.UUID, req dto.EndShiftRequest) (*dto.ShiftResponse, error)
	GetOpenShift(ctx context.Context, userID uuid.UUID) (*dto.ShiftResponse, error)
	CreateCashTransaction(ctx context.Context, userID uuid.UUID, req dto.CashTransactionRequest) (*dto.CashTransactionResponse, error)
}

type service struct {
	repo  repository.Store
	log   logger.ILogger
	cache *Cache
}

func NewService(repo repository.Store, log logger.ILogger, cache *Cache) Service {
	return &service{
		repo:  repo,
		log:   log,
		cache: cache,
	}
}

func (s *service) StartShift(ctx context.Context, userID uuid.UUID, req dto.StartShiftRequest) (*dto.ShiftResponse, error) {
	// Verify user password
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Errorf("StartShift | User not found: %v", userID)
		return nil, common.ErrNotFound
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
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

	return s.mapShiftToResponse(shift), nil
}

func (s *service) EndShift(ctx context.Context, userID uuid.UUID, req dto.EndShiftRequest) (*dto.ShiftResponse, error) {
	// Verify user password
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Errorf("EndShift | User not found: %v", userID)
		return nil, common.ErrNotFound
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
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

	expectedCashEnd := shift.StartCash + cashIn - cashOut
	// Note: We need to add Order Sales (Cash) here.
	// Assuming for now simple Shift management without strict order reconciliation in this step as it requires Order Repository dependency.

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
	s.cache.Clear(userID)

	return res, nil
}

func (s *service) GetOpenShift(ctx context.Context, userID uuid.UUID) (*dto.ShiftResponse, error) {
	shift, err := s.repo.GetOpenShiftByUserID(ctx, userID)
	if err != nil {
		return nil, common.ErrNotFound
	}

	return s.mapShiftToResponse(shift), nil
}

func (s *service) CreateCashTransaction(ctx context.Context, userID uuid.UUID, req dto.CashTransactionRequest) (*dto.CashTransactionResponse, error) {
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

	return &dto.CashTransactionResponse{
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

func (s *service) mapShiftToResponse(shift repository.Shift) *dto.ShiftResponse {
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

	return &dto.ShiftResponse{
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
