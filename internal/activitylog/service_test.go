package activitylog

import (
	"POS-kasir/internal/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
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
