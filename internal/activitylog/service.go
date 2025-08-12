package activitylog

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type IActivityService interface {
	Log(ctx context.Context, userID uuid.UUID, action repository.LogActionType, entityType repository.LogEntityType, entityID string, details map[string]interface{})
}

type ActivityService struct {
	repo repository.Store
	log  logger.ILogger
}

func NewActivityService(repo repository.Store, log logger.ILogger) IActivityService {
	return &ActivityService{
		repo: repo,
		log:  log,
	}
}

func (s *ActivityService) Log(ctx context.Context, userID uuid.UUID, action repository.LogActionType, entityType repository.LogEntityType, entityID string, details map[string]interface{}) {
	var detailsJSON []byte
	var err error

	if details != nil {
		detailsJSON, err = json.Marshal(details)
		if err != nil {
			s.log.Errorf("Log | Failed to marshal activity log details: %v", err)
		}
	}

	go func() {
		_, err := s.repo.CreateActivityLog(ctx, repository.CreateActivityLogParams{
			UserID:     pgtype.UUID{Bytes: userID, Valid: true},
			ActionType: action,
			EntityType: entityType,
			EntityID:   entityID,
			Details:    detailsJSON,
		})
		if err != nil {
			s.log.Errorf("Log | Failed to create activity log: %v", err)
		}
	}()
}
