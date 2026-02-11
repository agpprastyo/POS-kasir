package activitylog

import (
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestActivityService_Log(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	service := NewActivityService(mockRepo, mockLogger)
	ctx := context.Background()
	userID := uuid.New()
	entityID := "123"
	action := repository.LogActionTypeCREATE
	entityType := repository.LogEntityTypePRODUCT
	details := map[string]interface{}{"foo": "bar"}

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().CreateActivityLog(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, arg repository.CreateActivityLogParams) (uuid.UUID, error) {
				if arg.ActionType != action {
					t.Errorf("expected action %v, got %v", action, arg.ActionType)
				}
				if arg.EntityID != entityID {
					t.Errorf("expected entityID %v, got %v", entityID, arg.EntityID)
				}

				if string(arg.Details) != `{"foo":"bar"}` {
					t.Errorf("expected details JSON, got %s", arg.Details)
				}
				return uuid.New(), nil
			},
		)

		service.Log(ctx, userID, action, entityType, entityID, details)

		time.Sleep(50 * time.Millisecond)
	})

	t.Run("RepoError", func(t *testing.T) {

		mockRepo.EXPECT().CreateActivityLog(ctx, gomock.Any()).Return(uuid.Nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		service.Log(ctx, userID, action, entityType, entityID, details)

		time.Sleep(50 * time.Millisecond)
	})

	t.Run("JSONMarshalError", func(t *testing.T) {

		badDetails := map[string]interface{}{
			"bad": make(chan int),
		}

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		mockRepo.EXPECT().CreateActivityLog(ctx, gomock.Any()).Return(uuid.New(), nil)

		service.Log(ctx, userID, action, entityType, entityID, badDetails)

		time.Sleep(50 * time.Millisecond)
	})
}

func TestActivityService_GetActivityLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStore(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	service := NewActivityService(mockRepo, mockLogger)
	ctx := context.Background()

	t.Run("SuccessWithFilters", func(t *testing.T) {
		req := dto.GetActivityLogsRequest{
			Page:      1,
			Limit:     10,
			UserID:    uuid.New().String(),
			StartDate: "2023-01-01",
			EndDate:   "2023-01-02",
		}

		repoLogs := []repository.GetActivityLogsRow{
			{
				ID:         uuid.New(),
				UserID:     pgtype.UUID{Bytes: uuid.MustParse(req.UserID), Valid: true},
				ActionType: repository.LogActionTypeCREATE,
				EntityType: repository.LogEntityTypePRODUCT,
				EntityID:   "123",
				Details:    []byte(`{"foo":"bar"}`),
				CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
		}

		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return(repoLogs, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(1), nil)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Logs, 1)
		assert.Equal(t, int64(1), resp.TotalItems)
	})

	t.Run("SuccessNoFilters", func(t *testing.T) {
		req := dto.GetActivityLogsRequest{
			Page:  1,
			Limit: 10,
		}

		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return([]repository.GetActivityLogsRow{}, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(0), nil)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), resp.TotalItems)
	})

	t.Run("InvalidFilters", func(t *testing.T) {
		req := dto.GetActivityLogsRequest{
			Page:      1,
			Limit:     10,
			UserID:    "invalid-uuid",
			StartDate: "invalid-date",
			EndDate:   "invalid-date",
		}

		// It should proceed with empty filters
		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return([]repository.GetActivityLogsRow{}, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(0), nil)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("RepoError_GetActivityLogs", func(t *testing.T) {
		req := dto.GetActivityLogsRequest{Page: 1, Limit: 10}
		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("RepoError_CountActivityLogs", func(t *testing.T) {
		req := dto.GetActivityLogsRequest{Page: 1, Limit: 10}
		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return([]repository.GetActivityLogsRow{}, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(0), errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
