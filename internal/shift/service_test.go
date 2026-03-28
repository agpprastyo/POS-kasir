package shift_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/shift"
	repository "POS-kasir/internal/shift/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShiftService_StartShift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	mockCacheBase := mocks.NewMockCache(ctrl)
	shiftCache := shift.NewCache(mockCacheBase)
	service := shift.NewService(mockRepo, mockLogger, shiftCache)

	ctx := context.Background()
	userID := uuid.New()
	password := "password123"
	hash, _ := utils.HashPassword(password)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{}, errors.New("not found"))
		mockRepo.EXPECT().CreateShift(ctx, gomock.Any()).Return(repository.Shift{
			ID:        uuid.New(),
			UserID:    userID,
			StartCash: 100000,
			Status:    repository.ShiftStatus("open"),
		}, nil)
		mockCacheBase.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		resp, err := service.StartShift(ctx, userID, shift.StartShiftRequest{
			Password:  password,
			StartCash: 100000,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(100000), resp.StartCash)
		assert.Equal(t, repository.ShiftStatus("open"), resp.Status)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return("", errors.New("not found"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		resp, err := service.StartShift(ctx, userID, shift.StartShiftRequest{Password: password})
		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())

		resp, err := service.StartShift(ctx, userID, shift.StartShiftRequest{
			Password: "wrong",
		})

		assert.ErrorIs(t, err, common.ErrInvalidCredentials)
		assert.Nil(t, resp)
	})

	t.Run("AlreadyOpen", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())

		resp, err := service.StartShift(ctx, userID, shift.StartShiftRequest{
			Password: password,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already has an open shift")
		assert.Nil(t, resp)
	})

	t.Run("CreateError", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{}, errors.New("not found"))
		mockRepo.EXPECT().CreateShift(ctx, gomock.Any()).Return(repository.Shift{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.StartShift(ctx, userID, shift.StartShiftRequest{Password: password})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestShiftService_EndShift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	mockCacheBase := mocks.NewMockCache(ctrl)
	shiftCache := shift.NewCache(mockCacheBase)
	service := shift.NewService(mockRepo, mockLogger, shiftCache)

	ctx := context.Background()
	userID := uuid.New()
	password := "password123"
	hash, _ := utils.HashPassword(password)
	shiftID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{
			ID:        shiftID,
			StartCash: 100000,
		}, nil)
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(50000), nil).Times(1)
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(10000), nil).Times(1)
		
		expectedCashEnd := int64(140000)
		actualCashEnd := int64(135000)
		
		mockRepo.EXPECT().EndShift(ctx, gomock.Any()).Return(repository.Shift{
			ID:              shiftID,
			ExpectedCashEnd: &expectedCashEnd,
			ActualCashEnd:   &actualCashEnd,
			Status:          repository.ShiftStatus("closed"),
		}, nil)
		mockCacheBase.EXPECT().Delete(gomock.Any()).Return(nil)

		resp, err := service.EndShift(ctx, userID, shift.EndShiftRequest{
			Password:      password,
			ActualCashEnd: actualCashEnd,
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(-5000), *resp.Difference)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return("", errors.New("not found"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		_, err := service.EndShift(ctx, userID, shift.EndShiftRequest{Password: password})
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		_, err := service.EndShift(ctx, userID, shift.EndShiftRequest{Password: "wrong"})
		assert.ErrorIs(t, err, common.ErrInvalidCredentials)
	})

	t.Run("NoOpenShift", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{}, errors.New("not found"))
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())

		resp, err := service.EndShift(ctx, userID, shift.EndShiftRequest{
			Password: password,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no open shift found")
		assert.Nil(t, resp)
	})

	t.Run("CashInError", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{ID: shiftID}, nil)
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(0), errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		_, err := service.EndShift(ctx, userID, shift.EndShiftRequest{Password: password})
		assert.Error(t, err)
	})

	t.Run("CashOutError", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{ID: shiftID}, nil)
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(10), nil)
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(0), errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		_, err := service.EndShift(ctx, userID, shift.EndShiftRequest{Password: password})
		assert.Error(t, err)
	})

	t.Run("EndShiftRepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetUserPasswordHash(ctx, userID).Return(hash, nil)
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{ID: shiftID}, nil)
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(0), nil).Times(2)
		mockRepo.EXPECT().EndShift(ctx, gomock.Any()).Return(repository.Shift{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		_, err := service.EndShift(ctx, userID, shift.EndShiftRequest{Password: password})
		assert.Error(t, err)
	})
}

func TestShiftService_GetOpenShift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuerier(ctrl)
	service := shift.NewService(mockRepo, nil, nil)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{ID: uuid.New()}, nil)
		resp, err := service.GetOpenShift(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{}, errors.New("not found"))
		resp, err := service.GetOpenShift(ctx, userID)
		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})
}

func TestShiftService_CreateCashTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := shift.NewService(mockRepo, mockLogger, nil)

	ctx := context.Background()
	userID := uuid.New()
	shiftID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{ID: shiftID}, nil)
		mockRepo.EXPECT().CreateCashTransaction(ctx, gomock.Any()).Return(repository.CashTransaction{
			ID:     uuid.New(),
			Amount: 5000,
		}, nil)

		resp, err := service.CreateCashTransaction(ctx, userID, shift.CashTransactionRequest{
			Amount: 5000,
			Type:   repository.CashTransactionTypeCashIn,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(5000), resp.Amount)
	})

	t.Run("NoOpenShift", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{}, errors.New("not found"))
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any())
		_, err := service.CreateCashTransaction(ctx, userID, shift.CashTransactionRequest{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no open shift found")
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShiftByUserID(ctx, userID).Return(repository.Shift{ID: shiftID}, nil)
		mockRepo.EXPECT().CreateCashTransaction(ctx, gomock.Any()).Return(repository.CashTransaction{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		_, err := service.CreateCashTransaction(ctx, userID, shift.CashTransactionRequest{})
		assert.Error(t, err)
	})
}

func TestShiftService_AutoCloseShifts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	mockCacheBase := mocks.NewMockCache(ctrl)
	shiftCache := shift.NewCache(mockCacheBase)
	service := shift.NewService(mockRepo, mockLogger, shiftCache)

	ctx := context.Background()
	shiftID := uuid.New()
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShifts(ctx).Return([]repository.Shift{
			{ID: shiftID, UserID: userID, StartCash: 10000},
		}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any())
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(0), nil).AnyTimes()
		mockRepo.EXPECT().EndShift(ctx, gomock.Any()).Return(repository.Shift{}, nil)
		mockCacheBase.EXPECT().Delete(gomock.Any()).Return(nil)

		err := service.AutoCloseShifts(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetOpenShiftsError", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShifts(ctx).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any())
		err := service.AutoCloseShifts(ctx)
		assert.Error(t, err)
	})

	t.Run("PartialFailure", func(t *testing.T) {
		mockRepo.EXPECT().GetOpenShifts(ctx).Return([]repository.Shift{
			{ID: shiftID, UserID: userID},
			{ID: uuid.New(), UserID: uuid.New()},
		}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().GetCashTotalByShiftIDAndType(ctx, gomock.Any()).Return(int64(0), nil).AnyTimes()
		
		// First fails, second succeeds
		mockRepo.EXPECT().EndShift(ctx, gomock.Any()).Return(repository.Shift{}, errors.New("error")).Times(1)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		
		mockRepo.EXPECT().EndShift(ctx, gomock.Any()).Return(repository.Shift{}, nil).Times(1)
		mockCacheBase.EXPECT().Delete(gomock.Any()).Return(nil)

		err := service.AutoCloseShifts(ctx)
		assert.NoError(t, err) // Should continue despite partial failure
	})
}
