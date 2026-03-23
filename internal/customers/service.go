package customers

import (
	"context"
	"errors"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/customers/repository"
	"POS-kasir/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ICustomerService interface {
	CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*CustomerResponse, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (*CustomerResponse, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req UpdateCustomerRequest) (*CustomerResponse, error)
	DeleteCustomer(ctx context.Context, id uuid.UUID) error
	ListCustomers(ctx context.Context, req ListCustomersRequest) (*PagedCustomerResponse, error)
}

type CustomerService struct {
	repo repository.Querier
	log  logger.ILogger
}

func NewCustomerService(repo repository.Querier, log logger.ILogger) ICustomerService {
	return &CustomerService{repo: repo, log: log}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*CustomerResponse, error) {
	cust, err := s.repo.CreateCustomer(ctx, repository.CreateCustomerParams{
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	})
	if err != nil {
		s.log.Errorf("CreateCustomer failed", "error", err)
		return nil, err
	}
	return mapToCustomerResponse(cust), nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, id uuid.UUID) (*CustomerResponse, error) {
	cust, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return mapToCustomerResponse(cust), nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, req UpdateCustomerRequest) (*CustomerResponse, error) {
	cust, err := s.repo.UpdateCustomer(ctx, repository.UpdateCustomerParams{
		ID:      id,
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		s.log.Errorf("UpdateCustomer failed", "error", err)
		return nil, err
	}
	return mapToCustomerResponse(cust), nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteCustomer(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return common.ErrNotFound
		}
		s.log.Errorf("DeleteCustomer failed", "error", err)
		return err
	}
	return nil
}

func (s *CustomerService) ListCustomers(ctx context.Context, req ListCustomersRequest) (*PagedCustomerResponse, error) {
	req.SetDefaults()
	limit := req.Limit
	offset := (req.Page - 1) * limit

	custs, err := s.repo.ListCustomers(ctx, repository.ListCustomersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		s.log.Errorf("ListCustomers failed", "error", err)
		return nil, err
	}

	count, err := s.repo.CountCustomers(ctx)
	if err != nil {
		s.log.Errorf("CountCustomers failed", "error", err)
		return nil, err
	}

	var responses []CustomerResponse
	for _, c := range custs {
		responses = append(responses, *mapToCustomerResponse(c))
	}

	return &PagedCustomerResponse{
		Customers: responses,
		Pagination: pagination.BuildPagination(req.Page, int(count), limit),
	}, nil
}

func mapToCustomerResponse(c repository.Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Phone:     c.Phone,
		Email:     c.Email,
		Address:   c.Address,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: c.UpdatedAt.Time,
	}
}
