package products_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/products"
	"POS-kasir/mocks"
	"POS-kasir/pkg/utils"
	"POS-kasir/pkg/validator"
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPrdHandler_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPrdService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := products.NewPrdHandler(mockService, mockLogger)

	app := fiber.New(fiber.Config{
		StructValidator: validator.NewValidator(),
	})
	app.Use(recover.New())

	// Register all routes
	app.Get("/products/trash", handler.ListDeletedProductsHandler)
	app.Get("/products/trash/:id", handler.GetDeletedProductHandler)
	app.Post("/products/trash/:id/restore", handler.RestoreProductHandler)
	app.Post("/products/trash/restore", handler.RestoreProductsBulkHandler)
	
	app.Get("/products/:id", handler.GetProductHandler)
	app.Delete("/products/:id", handler.DeleteProductHandler)
	app.Post("/products", handler.CreateProductHandler)
	app.Get("/products", handler.ListProductsHandler)
	app.Patch("/products/:id", handler.UpdateProductHandler)
	app.Post("/products/:id/image", handler.UploadProductImageHandler)
	app.Get("/products/:id/stock-history", handler.GetStockHistoryHandler)
	app.Post("/products/:product_id/options", handler.CreateProductOptionHandler)
	app.Patch("/products/:product_id/options/:option_id", handler.UpdateProductOptionHandler)
	app.Post("/products/:product_id/options/:option_id/image", handler.UploadProductOptionImageHandler)
	app.Delete("/products/:product_id/options/:option_id", handler.DeleteProductOptionHandler)

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		req := httptest.NewRequest(http.MethodGet, "/products/invalid-uuid", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CreateValidationFailed", func(t *testing.T) {
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		reqBody := products.CreateProductRequest{Name: "sh"} // too short
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("GetNotFound", func(t *testing.T) {
		mockService.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(nil, common.ErrNotFound)
		req := httptest.NewRequest(http.MethodGet, "/products/"+uuid.New().String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		mockService.EXPECT().DeleteProduct(gomock.Any(), gomock.Any()).Return(common.ErrNotFound)
		req := httptest.NewRequest(http.MethodDelete, "/products/"+uuid.New().String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("RestoreNotFound", func(t *testing.T) {
		mockService.EXPECT().RestoreProduct(gomock.Any(), gomock.Any()).Return(common.ErrNotFound)
		req := httptest.NewRequest(http.MethodPost, "/products/trash/"+uuid.New().String()+"/restore", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("UpdateInternalError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().UpdateProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))
		body, _ := json.Marshal(products.UpdateProductRequest{})
		req := httptest.NewRequest(http.MethodPatch, "/products/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("RestoreBulkInvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products/trash/restore", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CreateInvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("{invalid}")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UploadNoFile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/image", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return(&products.ListProductsResponse{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/products?search=test&category_id=1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("UpdateOptionValidationFailed", func(t *testing.T) {
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		reqBody := products.UpdateProductOptionRequest{Name: utils.StringPtr("")} // too short, min=1
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPatch, "/products/"+uuid.New().String()+"/options/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UploadOptionNoFile", func(t *testing.T) {
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/options/"+uuid.New().String()+"/image", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("StockHistoryInternalError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().GetStockHistory(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))
		req := httptest.NewRequest(http.MethodGet, "/products/"+uuid.New().String()+"/stock-history", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("ListDeletedSuccess", func(t *testing.T) {
		mockService.EXPECT().ListDeletedProducts(gomock.Any(), gomock.Any()).Return(&products.ListProductsResponse{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/products/trash?page=1&limit=10", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetDeletedNotFound", func(t *testing.T) {
		mockService.EXPECT().GetDeletedProduct(gomock.Any(), gomock.Any()).Return(nil, common.ErrNotFound)
		req := httptest.NewRequest(http.MethodGet, "/products/trash/"+uuid.New().String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("UploadOptionServiceError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().UploadProductOptionImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))
		
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, _ := writer.CreateFormFile("image", "test.jpg")
		part.Write([]byte("data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/options/"+uuid.New().String()+"/image", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)
		if resp != nil {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("CreateOptionValidationFailed", func(t *testing.T) {
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		reqBody := products.CreateProductOptionRequestStandalone{Name: ""} // min=1
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/options", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CreateOptionServiceError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().CreateProductOption(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))
		body, _ := json.Marshal(products.CreateProductOptionRequestStandalone{Name: "valid", AdditionalPrice: 100})
		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/options", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("DeleteOptionNotFound", func(t *testing.T) {
		mockService.EXPECT().DeleteProductOption(gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrNotFound)
		req := httptest.NewRequest(http.MethodDelete, "/products/"+uuid.New().String()+"/options/"+uuid.New().String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("RestoreBulkServiceError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().RestoreProductsBulk(gomock.Any(), gomock.Any()).Return(errors.New("err"))
		body, _ := json.Marshal(products.RestoreBulkRequest{ProductIDs: []string{uuid.New().String()}})
		req := httptest.NewRequest(http.MethodPost, "/products/trash/restore", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("ListDeletedInternalError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().ListDeletedProducts(gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))
		req := httptest.NewRequest(http.MethodGet, "/products/trash", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetDeletedInternalError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().GetDeletedProduct(gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))
		productID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/products/trash/"+productID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("RestoreInternalError", func(t *testing.T) {
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockService.EXPECT().RestoreProduct(gomock.Any(), gomock.Any()).Return(errors.New("err"))
		productID := uuid.New()
		req := httptest.NewRequest(http.MethodPost, "/products/trash/"+productID.String()+"/restore", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetStockHistoryInvalidQuery", func(t *testing.T) {
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		productID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String()+"/stock-history?limit=invalid", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestPrdHandler_SuccessCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPrdService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	handler := products.NewPrdHandler(mockService, mockLogger)

	app := fiber.New()

	app.Get("/products/trash", handler.ListDeletedProductsHandler)
	app.Get("/products/trash/:id", handler.GetDeletedProductHandler)
	app.Post("/products/trash/:id/restore", handler.RestoreProductHandler)
	app.Post("/products/trash/restore", handler.RestoreProductsBulkHandler)
	app.Get("/products/:id/stock-history", handler.GetStockHistoryHandler)
	app.Post("/products/:id/image", handler.UploadProductImageHandler)
	app.Post("/products/:product_id/options/:option_id/image", handler.UploadProductOptionImageHandler)
	app.Delete("/products/:product_id/options/:option_id", handler.DeleteProductOptionHandler)
	app.Patch("/products/:product_id/options/:option_id", handler.UpdateProductOptionHandler)
	app.Post("/products/:product_id/options", handler.CreateProductOptionHandler)

	t.Run("UploadProductImage_Success", func(t *testing.T) {
		mockService.EXPECT().UploadProductImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(&products.ProductResponse{}, nil)
		
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, _ := writer.CreateFormFile("image", "test.jpg")
		part.Write([]byte("data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/image", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("UploadOptionImage_Success", func(t *testing.T) {
		mockService.EXPECT().UploadProductOptionImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&products.ProductOptionResponse{}, nil)
		
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, _ := writer.CreateFormFile("image", "test.jpg")
		part.Write([]byte("data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/options/"+uuid.New().String()+"/image", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DeleteOption_Success", func(t *testing.T) {
		mockService.EXPECT().DeleteProductOption(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		req := httptest.NewRequest(http.MethodDelete, "/products/"+uuid.New().String()+"/options/"+uuid.New().String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("UpdateOption_Success", func(t *testing.T) {
		mockService.EXPECT().UpdateProductOption(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&products.ProductOptionResponse{}, nil)
		body, _ := json.Marshal(products.UpdateProductOptionRequest{Name: utils.StringPtr("new name")})
		req := httptest.NewRequest(http.MethodPatch, "/products/"+uuid.New().String()+"/options/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("CreateOption_Success", func(t *testing.T) {
		mockService.EXPECT().CreateProductOption(gomock.Any(), gomock.Any(), gomock.Any()).Return(&products.ProductOptionResponse{}, nil)
		body, _ := json.Marshal(products.CreateProductOptionRequestStandalone{Name: "opt", AdditionalPrice: 100})
		req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/options", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("ListDeleted_Success", func(t *testing.T) {
		mockService.EXPECT().ListDeletedProducts(gomock.Any(), gomock.Any()).Return(&products.ListProductsResponse{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/products/trash", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetDeleted_Success", func(t *testing.T) {
		productID := uuid.New()
		mockService.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(&products.ProductResponse{ID: productID}, nil)
		req := httptest.NewRequest(http.MethodGet, "/products/trash/"+productID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("RestoreBulk_Success", func(t *testing.T) {
		mockService.EXPECT().RestoreProductsBulk(gomock.Any(), gomock.Any()).Return(nil)
		body, _ := json.Marshal(products.RestoreBulkRequest{ProductIDs: []string{uuid.New().String()}})
		req := httptest.NewRequest(http.MethodPost, "/products/trash/restore", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetStockHistory_Success", func(t *testing.T) {
		productID := uuid.New()
		mockService.EXPECT().GetStockHistory(gomock.Any(), productID, gomock.Any()).Return(&products.PagedStockHistoryResponse{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String()+"/stock-history", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
