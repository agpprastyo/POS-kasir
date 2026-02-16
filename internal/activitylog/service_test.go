package activitylog_test

import (
	"POS-kasir/internal/activitylog"
	activitylog_repo "POS-kasir/internal/activitylog/repository"
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

	mockRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	service := activitylog.NewActivityService(mockRepo, mockLogger)
	ctx := context.Background()
	userID := uuid.New()
	entityID := "123"
	action := activitylog_repo.LogActionTypeCREATE
	entityType := activitylog_repo.LogEntityTypePRODUCT
	details := map[string]interface{}{"foo": "bar"}

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().CreateActivityLog(ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, arg activitylog_repo.CreateActivityLogParams) (uuid.UUID, error) {
				if arg.ActionType != activitylog_repo.LogActionType(action) {
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

	mockRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	service := activitylog.NewActivityService(mockRepo, mockLogger)
	ctx := context.Background()

	t.Run("SuccessWithFilters", func(t *testing.T) {
		page, limit := 1, 10
		userID, startDate, endDate := uuid.New().String(), "2023-01-01", "2023-01-02"
		req := activitylog.GetActivityLogsRequest{
			Page:      &page,
			Limit:     &limit,
			UserID:    &userID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		repoLogs := []activitylog_repo.GetActivityLogsRow{
			{
				ID:         uuid.New(),
				UserID:     pgtype.UUID{Bytes: uuid.MustParse(*req.UserID), Valid: true},
				ActionType: activitylog_repo.LogActionTypeCREATE,
				EntityType: activitylog_repo.LogEntityTypePRODUCT,
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
		page, limit := 1, 10
		req := activitylog.GetActivityLogsRequest{
			Page:  &page,
			Limit: &limit,
		}

		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return([]activitylog_repo.GetActivityLogsRow{}, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(0), nil)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), resp.TotalItems)
	})

	t.Run("InvalidFilters", func(t *testing.T) {
		page, limit := 1, 10
		userID, startDate, endDate := "invalid-uuid", "invalid-date", "invalid-date"
		req := activitylog.GetActivityLogsRequest{
			Page:      &page,
			Limit:     &limit,
			UserID:    &userID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		// It should proceed with empty filters
		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return([]activitylog_repo.GetActivityLogsRow{}, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(0), nil)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("RepoError_GetActivityLogs", func(t *testing.T) {
		page, limit := 1, 10
		req := activitylog.GetActivityLogsRequest{Page: &page, Limit: &limit}
		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("RepoError_CountActivityLogs", func(t *testing.T) {
		page, limit := 1, 10
		req := activitylog.GetActivityLogsRequest{Page: &page, Limit: &limit}
		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return([]activitylog_repo.GetActivityLogsRow{}, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(0), errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("UnmarshalFailure_Details", func(t *testing.T) {
		page, limit := 1, 10
		req := activitylog.GetActivityLogsRequest{Page: &page, Limit: &limit}

		repoLogs := []activitylog_repo.GetActivityLogsRow{
			{
				ID:         uuid.New(),
				ActionType: activitylog_repo.LogActionTypeCREATE,
				EntityType: activitylog_repo.LogEntityTypePRODUCT,
				Details:    []byte(`{invalid json}`),
				CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
		}

		mockRepo.EXPECT().GetActivityLogs(ctx, gomock.Any()).Return(repoLogs, nil)
		mockRepo.EXPECT().CountActivityLogs(ctx, gomock.Any()).Return(int64(1), nil)

		resp, err := service.GetActivityLogs(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Logs, 1)
		// Details should be nil if unmarshal fails but error is ignored in code
		assert.Nil(t, resp.Logs[0].Details)
	})
}
