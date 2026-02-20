package promotions_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/promotions"
	promo_repo "POS-kasir/internal/promotions/repository"
	"POS-kasir/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPromotionHandler_CreatePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPromotionService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	handler := promotions.NewPromotionHandler(mockService, mockLogger)
	app := fiber.New()
	app.Post("/promotions", handler.CreatePromotionHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := promotions.CreatePromotionRequest{
			Name:          "Promo Merdeka",
			Scope:         promo_repo.PromotionScopeORDER,
			DiscountType:  promo_repo.DiscountTypePercentage,
			DiscountValue: 10,
			StartDate:     time.Now(),
			EndDate:       time.Now().Add(24 * time.Hour),
			IsActive:      true,
			Rules:         []promotions.CreatePromotionRuleRequest{},
			Targets:       []promotions.CreatePromotionTargetRequest{},
		}

		expectedResp := &promotions.PromotionResponse{
			ID:   uuid.New(),
			Name: reqBody.Name,
		}

		mockService.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Return(expectedResp, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/promotions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/promotions", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := promotions.CreatePromotionRequest{
			Name: "Promo Error",
		}

		mockService.EXPECT().CreatePromotion(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/promotions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		}
	})
}

func TestPromotionHandler_UpdatePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPromotionService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	handler := promotions.NewPromotionHandler(mockService, mockLogger)
	app := fiber.New()
	app.Put("/promotions/:id", handler.UpdatePromotionHandler)

	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		reqBody := promotions.UpdatePromotionRequest{
			Name: "Promo Update",
		}

		expectedResp := &promotions.PromotionResponse{
			ID:   promoID,
			Name: reqBody.Name,
		}

		mockService.EXPECT().UpdatePromotion(gomock.Any(), promoID, gomock.Any()).Return(expectedResp, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/promotions/"+promoID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestPromotionHandler_GetPromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPromotionService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	handler := promotions.NewPromotionHandler(mockService, mockLogger)
	app := fiber.New()
	app.Get("/promotions/:id", handler.GetPromotionHandler)

	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedResp := &promotions.PromotionResponse{
			ID:   promoID,
			Name: "Promo Get",
		}

		mockService.EXPECT().GetPromotion(gomock.Any(), promoID).Return(expectedResp, nil)

		req := httptest.NewRequest(http.MethodGet, "/promotions/"+promoID.String(), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().GetPromotion(gomock.Any(), promoID).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/promotions/"+promoID.String(), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		}
	})
}

func TestPromotionHandler_ListPromotions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPromotionService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	handler := promotions.NewPromotionHandler(mockService, mockLogger)
	app := fiber.New()
	app.Get("/promotions", handler.ListPromotionsHandler)

	t.Run("Success", func(t *testing.T) {
		expectedResp := &promotions.PagedPromotionResponse{
			Promotions: []promotions.PromotionResponse{
				{ID: uuid.New(), Name: "Promo List"},
			},
			Pagination: pagination.Pagination{
				CurrentPage: 1,
				PerPage:     10,
				TotalData:   1,
				TotalPage:   1,
			},
		}

		mockService.EXPECT().ListPromotions(gomock.Any(), gomock.Any()).Return(expectedResp, nil)

		req := httptest.NewRequest(http.MethodGet, "/promotions?page=1&limit=10", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})
}

func TestPromotionHandler_DeletePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPromotionService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	handler := promotions.NewPromotionHandler(mockService, mockLogger)
	app := fiber.New()
	app.Delete("/promotions/:id", handler.DeletePromotionHandler)

	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().DeletePromotion(gomock.Any(), promoID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/promotions/"+promoID.String(), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})
}

func TestPromotionHandler_RestorePromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPromotionService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	handler := promotions.NewPromotionHandler(mockService, mockLogger)
	app := fiber.New()
	app.Post("/promotions/:id/restore", handler.RestorePromotionHandler)

	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().RestorePromotion(gomock.Any(), promoID).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/promotions/"+promoID.String()+"/restore", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})
}
