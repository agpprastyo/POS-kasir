package payment_methods

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
)

type IPaymentMethodService interface {
	ListPaymentMethods(ctx context.Context) ([]PaymentMethodResponse, error)
}

type PaymentMethodService struct {
	store repository.Store
	log   logger.ILogger
}

func NewPaymentMethodService(store repository.Store, log logger.ILogger) IPaymentMethodService {
	return &PaymentMethodService{store: store, log: log}
}

func (s *PaymentMethodService) ListPaymentMethods(ctx context.Context) ([]PaymentMethodResponse, error) {
	methods, err := s.store.ListPaymentMethods(ctx)
	if err != nil {
		s.log.Error("Failed to list payment methods from repository", "error", err)
		return nil, err
	}

	var response []PaymentMethodResponse
	for _, method := range methods {
		response = append(response, PaymentMethodResponse{
			ID:        method.ID,
			Name:      method.Name,
			IsActive:  method.IsActive,
			CreatedAt: method.CreatedAt.Time,
		})
	}
	return response, nil
}
