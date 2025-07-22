package activitylog

import (
	"POS-kasir/internal/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestLogService(t *testing.T) {
	// Generate a fixed UUID for consistent testing
	testUserID := uuid.New()

	testCases := []struct {
		name        string
		userID      uuid.UUID
		action      repository.LogActionType
		entityType  repository.LogEntityType
		entityID    string
		details     map[string]interface{}
		setupMocks  func(mockRepo *mocks.MockStore, mockLogger *mocks.MockILogger, wg *sync.WaitGroup)
		expectError bool
	}{
		{
			name:       "Success Case",
			userID:     testUserID,
			action:     repository.LogActionTypeCREATE,
			entityType: repository.LogEntityTypePRODUCT,
			entityID:   "prod_123",
			details:    map[string]interface{}{"name": "New Product"},
			setupMocks: func(mockRepo *mocks.MockStore, mockLogger *mocks.MockILogger, wg *sync.WaitGroup) {
				// We expect CreateActivityLog to be called once.
				mockRepo.EXPECT().
					CreateActivityLog(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params repository.CreateActivityLogParams) (uuid.UUID, error) {
						defer wg.Done() // Signal that this function has been called.
						// CORRECTED: Return a uuid.UUID as indicated by the error.
						return uuid.New(), nil
					})
			},
			expectError: false,
		},
		{
			name:       "Success Case with Nil Details",
			userID:     testUserID,
			action:     repository.LogActionTypeDELETE,
			entityType: repository.LogEntityTypeCATEGORY,
			entityID:   "cat_456",
			details:    nil,
			setupMocks: func(mockRepo *mocks.MockStore, mockLogger *mocks.MockILogger, wg *sync.WaitGroup) {
				// Expect CreateActivityLog to be called with nil details.
				mockRepo.EXPECT().
					CreateActivityLog(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params repository.CreateActivityLogParams) (uuid.UUID, error) {
						defer wg.Done()
						if params.Details != nil {
							t.Errorf("expected details to be nil, got %v", params.Details)
						}
						// CORRECTED: Return a uuid.UUID as indicated by the error.
						return uuid.New(), nil
					})
			},
			expectError: false,
		},
		{
			name:       "Repository Error Case",
			userID:     testUserID,
			action:     repository.LogActionTypeUPDATE,
			entityType: repository.LogEntityTypeUSER,
			entityID:   "user_789",
			details:    map[string]interface{}{"field": "password"},
			setupMocks: func(mockRepo *mocks.MockStore, mockLogger *mocks.MockILogger, wg *sync.WaitGroup) {
				// Mock the repository to return an error.
				// CORRECTED: Return a zero-value uuid.UUID and the error.
				mockRepo.EXPECT().
					CreateActivityLog(gomock.Any(), gomock.Any()).
					Return(uuid.UUID{}, errors.New("database connection failed")).
					Times(1) // It will be called once.

				// Expect the logger to be called with the error.
				mockLogger.EXPECT().
					Errorf(gomock.Eq("Failed to create activity log"), gomock.Any()).
					DoAndReturn(func(template string, args ...interface{}) {
						defer wg.Done() // Signal that the logger was called.
					})
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockStore(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)

			// Use a WaitGroup to synchronize with the goroutine in the Log method.
			var wg sync.WaitGroup
			wg.Add(1) // We expect one async operation to complete.

			// Setup the mock expectations for the current test case.
			tc.setupMocks(mockRepo, mockLogger, &wg)

			// Create the ActivityService with the mocks.
			s := NewActivityService(mockRepo, mockLogger)

			// Call the method under test.
			s.Log(context.Background(), tc.userID, tc.action, tc.entityType, tc.entityID, tc.details)

			// Wait for the goroutine to finish its work, with a timeout to prevent tests from hanging.
			waitTimeout(&wg, 1*time.Second, t)
		})
	}
}

// waitTimeout waits for the waitgroup for the specified duration.
// If the waitgroup does not finish in time, it fails the test.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration, t *testing.T) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		// Waitgroup finished as expected.
	case <-time.After(timeout):
		t.Fatal("Test timed out waiting for goroutine to finish")
	}
}
