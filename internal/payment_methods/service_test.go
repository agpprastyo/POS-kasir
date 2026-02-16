package payment_methods_test

import (
	"POS-kasir/internal/payment_methods"
	"POS-kasir/internal/payment_methods/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentMethodService_ListPaymentMethods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPaymentMethodsRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	service := payment_methods.NewPaymentMethodService(mockRepo, mockLogger)

	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		repoMethods := []repository.PaymentMethod{
			{
				ID:        1,
				Name:      "Cash",
				IsActive:  true,
				CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			},
			{
				ID:        2,
				Name:      "QRIS",
				IsActive:  true,
				CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			},
		}

		mockRepo.EXPECT().ListPaymentMethods(ctx).Return(repoMethods, nil)

		resp, err := service.ListPaymentMethods(ctx)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "Cash", resp[0].Name)
		assert.Equal(t, "QRIS", resp[1].Name)
		assert.True(t, resp[0].IsActive)
		assert.Equal(t, now, resp[0].CreatedAt)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		dbErr := errors.New("database error")
		mockRepo.EXPECT().ListPaymentMethods(ctx).Return(nil, dbErr)
		mockLogger.EXPECT().Error("Failed to list payment methods from repository", "error", dbErr)

		resp, err := service.ListPaymentMethods(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, dbErr, err)
	})
}
