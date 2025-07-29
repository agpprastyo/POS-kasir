package cancellation_reasons

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
)

type ICancellationReasonService interface {
	ListCancellationReasons(ctx context.Context) ([]CancellationReasonResponse, error)
}

type CancellationReasonService struct {
	store repository.Store
	log   logger.ILogger
}

func NewCancellationReasonService(store repository.Store, log logger.ILogger) ICancellationReasonService {
	return &CancellationReasonService{store: store, log: log}
}

func (s *CancellationReasonService) ListCancellationReasons(ctx context.Context) ([]CancellationReasonResponse, error) {
	reasons, err := s.store.ListCancellationReasons(ctx)
	if err != nil {
		s.log.Error("ListCancellationReasons | Failed to list cancellation reasons from repository", "error", err)
		return nil, err
	}

	var response []CancellationReasonResponse
	for _, reason := range reasons {
		response = append(response, CancellationReasonResponse{
			ID:          reason.ID,
			Reason:      reason.Reason,
			Description: reason.Description,
			IsActive:    reason.IsActive,
			CreatedAt:   reason.CreatedAt.Time,
		})
	}
	return response, nil
}
