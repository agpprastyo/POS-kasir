package cancellation_reasons

import (
	"POS-kasir/internal/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupServiceTest(t *testing.T) (*mocks.MockStore, *mocks.MockFieldLogger, ICancellationReasonService) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	service := NewCancellationReasonService(mockStore, mockLogger)
	return mockStore, mockLogger, service
}

func TestCancellationReasonService_ListCancellationReasons(t *testing.T) {
	mockStore, mockLogger, service := setupServiceTest(t)
	ctx := context.Background()
	now := time.Now()
	desc := "Default descriptions"

	t.Run("Success", func(t *testing.T) {
		repoReasons := []repository.CancellationReason{
			{
				ID:          1,
				Reason:      "Reason 1",
				Description: &desc,
				IsActive:    true,
				CreatedAt:   pgtype.Timestamptz{Time: now, Valid: true},
			},
		}

		mockStore.EXPECT().ListCancellationReasons(ctx).Return(repoReasons, nil)

		resp, err := service.ListCancellationReasons(ctx)

		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, int32(1), resp[0].ID)
		assert.Equal(t, "Reason 1", resp[0].Reason)
		assert.Equal(t, &desc, resp[0].Description)
		assert.True(t, resp[0].IsActive)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockStore.EXPECT().ListCancellationReasons(ctx).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.ListCancellationReasons(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "db error", err.Error())
	})
}
