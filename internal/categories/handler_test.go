package categories_test

import (
	"POS-kasir/internal/categories"
	"POS-kasir/internal/common"
	"POS-kasir/mocks"
	"POS-kasir/pkg/validator"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*mocks.MockICtgService, *mocks.MockFieldLogger, *categories.CtgHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockICtgService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	handler := categories.NewCtgHandler(mockService, mockLogger).(*categories.CtgHandler)
	app := fiber.New(fiber.Config{
		StructValidator: validator.NewValidator(),
	})
	return mockService, mockLogger, handler, app
}

func TestCtgHandler_GetAllCategoriesHandler(t *testing.T) {
	mockService, _, handler, app := setupHandlerTest(t)
	app.Get("/categories", handler.GetAllCategoriesHandler)

	t.Run("Success", func(t *testing.T) {
		categoriesList := []categories.CategoryResponse{
			{ID: 1, Name: "Category 1", CreatedAt: time.Now()},
		}
		mockService.EXPECT().GetAllCategories(gomock.Any(), gomock.Any()).Return(categoriesList, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().GetAllCategories(gomock.Any(), gomock.Any()).Return(nil, common.ErrCategoryNotFound)

		req := httptest.NewRequest("GET", "/categories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCtgHandler_CreateCategoryHandler(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Post("/categories", handler.CreateCategoryHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "New Category"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CreateCategory(gomock.Any(), reqBody).Return(&categories.CategoryResponse{ID: 1, Name: "New Category"}, nil)

		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Conflict", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "Existing"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CreateCategory(gomock.Any(), reqBody).Return(nil, common.ErrCategoryExists)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "Fail"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CreateCategory(gomock.Any(), reqBody).Return(nil, errors.New("error"))
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingName", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: ""}
		body, _ := json.Marshal(reqBody)

		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestCtgHandler_GetCategoryByIDHandler(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Get("/categories/:id", handler.GetCategoryByIDHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().GetCategoryByID(gomock.Any(), int32(1)).Return(&categories.CategoryResponse{ID: 1, Name: "Found"}, nil)

		req := httptest.NewRequest("GET", "/categories/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().GetCategoryByID(gomock.Any(), int32(1)).Return(nil, nil)

		req := httptest.NewRequest("GET", "/categories/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService.EXPECT().GetCategoryByID(gomock.Any(), int32(1)).Return(nil, errors.New("error"))
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/categories/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/categories/abc", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestCtgHandler_UpdateCategoryHandler(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Put("/categories/:id", handler.UpdateCategoryHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "Updated"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateCategory(gomock.Any(), int32(1), reqBody).Return(&categories.CategoryResponse{ID: 1, Name: "Updated"}, nil)

		req := httptest.NewRequest("PUT", "/categories/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "UpdateMe"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateCategory(gomock.Any(), int32(99), reqBody).Return(nil, common.ErrCategoryNotFound)

		req := httptest.NewRequest("PUT", "/categories/99", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Conflict", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "Exists"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateCategory(gomock.Any(), int32(1), reqBody).Return(nil, common.ErrCategoryExists)

		req := httptest.NewRequest("PUT", "/categories/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: "Error"}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateCategory(gomock.Any(), int32(1), reqBody).Return(nil, errors.New("error"))

		req := httptest.NewRequest("PUT", "/categories/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/categories/abc", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/categories/1", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingName", func(t *testing.T) {
		reqBody := categories.CreateCategoryRequest{Name: ""}
		body, _ := json.Marshal(reqBody)

		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/categories/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestCtgHandler_DeleteCategoryHandler(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Delete("/categories/:id", handler.DeleteCategoryHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().DeleteCategory(gomock.Any(), int32(1)).Return(nil)

		req := httptest.NewRequest("DELETE", "/categories/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().DeleteCategory(gomock.Any(), int32(99)).Return(common.ErrCategoryNotFound)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("DELETE", "/categories/99", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("ConflictInUse", func(t *testing.T) {
		mockService.EXPECT().DeleteCategory(gomock.Any(), int32(1)).Return(common.ErrCategoryInUse)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("DELETE", "/categories/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService.EXPECT().DeleteCategory(gomock.Any(), int32(1)).Return(errors.New("error"))
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("DELETE", "/categories/1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("DELETE", "/categories/abc", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestCtgHandler_GetCategoryCountHandler(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Get("/categories/count", handler.GetCategoryCountHandler)

	t.Run("Success", func(t *testing.T) {
		counts := []categories.CategoryWithCountResponse{{ID: 1, Name: "C1", ProductCount: 5}}
		mockService.EXPECT().GetCategoryWithProductCount(gomock.Any()).Return(&counts, nil)

		req := httptest.NewRequest("GET", "/categories/count", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Error", func(t *testing.T) {
		mockService.EXPECT().GetCategoryWithProductCount(gomock.Any()).Return(nil, errors.New("error"))
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/categories/count", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestCtgHandler_GetAllCategoriesHandler_Errors(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Get("/categories", handler.GetAllCategoriesHandler)

	t.Run("QueryParseError", func(t *testing.T) {
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/categories?limit=abc", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService.EXPECT().GetAllCategories(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

		req := httptest.NewRequest("GET", "/categories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		mockService.EXPECT().GetAllCategories(gomock.Any(), gomock.Any()).Return([]categories.CategoryResponse{}, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
