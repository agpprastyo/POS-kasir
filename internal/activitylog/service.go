// file: internal/activitylog/service.go

package activitylog

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	Log(ctx context.Context, userID uuid.UUID, action repository.LogActionType, entityType repository.LogEntityType, entityID string, details map[string]interface{})
}

type service struct {
	repo repository.Querier
	log  *logger.Logger
}

func NewService(repo repository.Querier, log *logger.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Log(ctx context.Context, userID uuid.UUID, action repository.LogActionType, entityType repository.LogEntityType, entityID string, details map[string]interface{}) {
	var detailsJSON []byte
	var err error

	if details != nil {
		detailsJSON, err = json.Marshal(details)
		if err != nil {
			s.log.Error("Failed to marshal activity log details", "error", err)
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
			s.log.Error("Failed to create activity log", "error", err)
		}
	}()
}
