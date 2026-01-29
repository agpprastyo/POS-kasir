package activitylog

import (
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type IActivityService interface {
	Log(ctx context.Context, userID uuid.UUID, action repository.LogActionType, entityType repository.LogEntityType, entityID string, details map[string]interface{})
	GetActivityLogs(ctx context.Context, req dto.GetActivityLogsRequest) (*dto.ActivityLogListResponse, error)
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

func (s *ActivityService) GetActivityLogs(ctx context.Context, req dto.GetActivityLogsRequest) (*dto.ActivityLogListResponse, error) {
	arg := repository.GetActivityLogsParams{
		Limit:  int32(req.Limit),
		Offset: int32((req.Page - 1) * req.Limit),
	}

	countArg := repository.CountActivityLogsParams{}

	if req.UserID != "" {
		uid, err := uuid.Parse(req.UserID)
		if err == nil {
			arg.UserID = pgtype.UUID{Bytes: uid, Valid: true}
			countArg.UserID = pgtype.UUID{Bytes: uid, Valid: true}
		}
	}

	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			arg.StartDate = pgtype.Timestamptz{Time: t, Valid: true}
			countArg.StartDate = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}

	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			// Set to end of day
			t = t.Add(24 * time.Hour).Add(-1 * time.Second)
			arg.EndDate = pgtype.Timestamptz{Time: t, Valid: true}
			countArg.EndDate = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}

	logs, err := s.repo.GetActivityLogs(ctx, arg)
	if err != nil {
		s.log.Errorf("GetActivityLogs | Failed to fetch logs: %v", err)
		return nil, err
	}

	totalItems, err := s.repo.CountActivityLogs(ctx, countArg)
	if err != nil {
		s.log.Errorf("GetActivityLogs | Failed to count logs: %v", err)
		return nil, err
	}

	response := &dto.ActivityLogListResponse{
		Logs:       make([]dto.ActivityLogResponse, 0),
		TotalItems: totalItems,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int(math.Ceil(float64(totalItems) / float64(req.Limit))),
	}

	for _, log := range logs {
		var details map[string]interface{}
		if len(log.Details) > 0 {
			_ = json.Unmarshal(log.Details, &details)
		}

		var userName string
		if log.UserName != nil {
			userName = *log.UserName
		}

		response.Logs = append(response.Logs, dto.ActivityLogResponse{
			ID:         log.ID,
			UserID:     log.UserID.Bytes,
			UserName:   userName,
			ActionType: log.ActionType,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			Details:    details,
			CreatedAt:  log.CreatedAt.Time,
		})
	}

	return response, nil
}
