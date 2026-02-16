package payment_methods

import (
	"POS-kasir/internal/payment_methods/repository"
	"POS-kasir/pkg/logger"
	"context"
)

type IPaymentMethodService interface {
	ListPaymentMethods(ctx context.Context) ([]PaymentMethodResponse, error)
}

type PaymentMethodService struct {
	repo repository.Querier
	log  logger.ILogger
}

func NewPaymentMethodService(repo repository.Querier, log logger.ILogger) IPaymentMethodService {
	return &PaymentMethodService{repo: repo, log: log}
}

func (s *PaymentMethodService) ListPaymentMethods(ctx context.Context) ([]PaymentMethodResponse, error) {
	methods, err := s.repo.ListPaymentMethods(ctx)
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
