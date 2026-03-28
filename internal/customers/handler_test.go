package customers_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/customers"
	"POS-kasir/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCustomerHandler_CreateCustomerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockICustomerService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := customers.NewCustomerHandler(mockService, mockLogger)

	app := fiber.New()
	app.Post("/customers", handler.CreateCustomerHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := customers.CreateCustomerRequest{Name: "John Doe"}
		respData := &customers.CustomerResponse{ID: uuid.New(), Name: "John Doe"}

		mockService.EXPECT().CreateCustomer(gomock.Any(), reqBody).Return(respData, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := customers.CreateCustomerRequest{Name: "John Doe"}
		mockService.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCustomerHandler_GetCustomerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockICustomerService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := customers.NewCustomerHandler(mockService, mockLogger)

	app := fiber.New()
	app.Get("/customers/:id", handler.GetCustomerHandler)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().GetCustomer(gomock.Any(), id).Return(&customers.CustomerResponse{ID: id}, nil)

		req := httptest.NewRequest(http.MethodGet, "/customers/"+id.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/customers/invalid", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().GetCustomer(gomock.Any(), id).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/customers/"+id.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCustomerHandler_UpdateCustomerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockICustomerService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := customers.NewCustomerHandler(mockService, mockLogger)

	app := fiber.New()
	app.Put("/customers/:id", handler.UpdateCustomerHandler)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		reqBody := customers.UpdateCustomerRequest{Name: "Jane Doe"}
		mockService.EXPECT().UpdateCustomer(gomock.Any(), id, reqBody).Return(&customers.CustomerResponse{ID: id, Name: "Jane Doe"}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/customers/"+id.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Error", func(t *testing.T) {
		reqBody := customers.UpdateCustomerRequest{Name: "Jane Doe"}
		mockService.EXPECT().UpdateCustomer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/customers/"+id.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCustomerHandler_DeleteCustomerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockICustomerService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := customers.NewCustomerHandler(mockService, mockLogger)

	app := fiber.New()
	app.Delete("/customers/:id", handler.DeleteCustomerHandler)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().DeleteCustomer(gomock.Any(), id).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/customers/"+id.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Error", func(t *testing.T) {
		mockService.EXPECT().DeleteCustomer(gomock.Any(), id).Return(errors.New("error"))

		req := httptest.NewRequest(http.MethodDelete, "/customers/"+id.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCustomerHandler_ListCustomersHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockICustomerService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := customers.NewCustomerHandler(mockService, mockLogger)

	app := fiber.New()
	app.Get("/customers", handler.ListCustomersHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().ListCustomers(gomock.Any(), gomock.Any()).Return(&customers.PagedCustomerResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Error", func(t *testing.T) {
		mockService.EXPECT().ListCustomers(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
