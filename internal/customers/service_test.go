package customers_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/customers"
	"POS-kasir/internal/customers/repository"
	"POS-kasir/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func strPtr(s string) *string {
	return &s
}

func TestCustomerService_CreateCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := customers.NewCustomerService(mockRepo, mockLogger)

	ctx := context.Background()
	req := customers.CreateCustomerRequest{
		Name:    "John Doe",
		Phone:   strPtr("08123456789"),
		Email:   strPtr("john@example.com"),
		Address: strPtr("Jl. Merdeka No. 1"),
	}

	customer := repository.Customer{
		ID:      uuid.New(),
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().CreateCustomer(ctx, repository.CreateCustomerParams{
			Name:    req.Name,
			Phone:   req.Phone,
			Email:   req.Email,
			Address: req.Address,
		}).Return(customer, nil)

		resp, err := service.CreateCustomer(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, customer.ID, resp.ID)
		assert.Equal(t, customer.Name, resp.Name)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().CreateCustomer(ctx, gomock.Any()).Return(repository.Customer{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.CreateCustomer(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestCustomerService_GetCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := customers.NewCustomerService(mockRepo, mockLogger)

	ctx := context.Background()
	id := uuid.New()
	customer := repository.Customer{ID: id, Name: "John Doe"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetCustomerByID(ctx, id).Return(customer, nil)

		resp, err := service.GetCustomer(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetCustomerByID(ctx, id).Return(repository.Customer{}, pgx.ErrNoRows)

		resp, err := service.GetCustomer(ctx, id)
		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetCustomerByID(ctx, id).Return(repository.Customer{}, errors.New("db error"))

		resp, err := service.GetCustomer(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestCustomerService_UpdateCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := customers.NewCustomerService(mockRepo, mockLogger)

	ctx := context.Background()
	id := uuid.New()
	req := customers.UpdateCustomerRequest{Name: "Jane Doe"}
	customer := repository.Customer{ID: id, Name: "Jane Doe"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().UpdateCustomer(ctx, repository.UpdateCustomerParams{
			ID:   id,
			Name: req.Name,
		}).Return(customer, nil)

		resp, err := service.UpdateCustomer(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, "Jane Doe", resp.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().UpdateCustomer(ctx, gomock.Any()).Return(repository.Customer{}, pgx.ErrNoRows)

		resp, err := service.UpdateCustomer(ctx, id, req)
		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().UpdateCustomer(ctx, gomock.Any()).Return(repository.Customer{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.UpdateCustomer(ctx, id, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestCustomerService_DeleteCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := customers.NewCustomerService(mockRepo, mockLogger)

	ctx := context.Background()
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().DeleteCustomer(ctx, id).Return(nil)

		err := service.DeleteCustomer(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().DeleteCustomer(ctx, id).Return(pgx.ErrNoRows)

		err := service.DeleteCustomer(ctx, id)
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteCustomer(ctx, id).Return(errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		err := service.DeleteCustomer(ctx, id)
		assert.Error(t, err)
	})
}

func TestCustomerService_ListCustomers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	service := customers.NewCustomerService(mockRepo, mockLogger)

	ctx := context.Background()
	req := customers.ListCustomersRequest{}
	req.SetDefaults()

	custs := []repository.Customer{{ID: uuid.New(), Name: "John Doe"}}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().ListCustomers(ctx, gomock.Any()).Return(custs, nil)
		mockRepo.EXPECT().CountCustomers(ctx).Return(int64(1), nil)

		resp, err := service.ListCustomers(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, resp.Customers, 1)
		assert.Equal(t, 1, resp.Pagination.TotalData)
	})

	t.Run("ListError", func(t *testing.T) {
		mockRepo.EXPECT().ListCustomers(ctx, gomock.Any()).Return(nil, errors.New("list error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.ListCustomers(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("CountError", func(t *testing.T) {
		mockRepo.EXPECT().ListCustomers(ctx, gomock.Any()).Return(custs, nil)
		mockRepo.EXPECT().CountCustomers(ctx).Return(int64(0), errors.New("count error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.ListCustomers(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
