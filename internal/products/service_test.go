package products_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/products"
	products_repo "POS-kasir/internal/products/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPrdService_GetProductByID(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	imageUrl := "products/test.jpg"

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockImageRepo := mocks.NewMockIProductImageRepository(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImageRepo, nil)

		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{
			ID:         productID,
			Name:       "Test Product",
			ImageUrl:   &imageUrl,
			Price:      10000,
			Stock:      10,
			CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
			CostPrice:  pgtype.Numeric{Valid: true},
			Categories: []map[string]interface{}{{"id": int32(1), "name": "Category 1"}},
			Options:    []map[string]interface{}{{"id": uuid.Nil, "name": "Option 1", "additional_price": int64(5000)}},
		}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		mockImageRepo.EXPECT().PrdImageLink(gomock.Any(), gomock.Any(), gomock.Any()).Return("http://public.url/test.jpg", nil).AnyTimes()

		resp, err := service.GetProductByID(ctx, productID)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestPrdService_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	newName := "Updated Name"
	req := products.UpdateProductRequest{
		Name: &newName,
	}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{ID: productID, Name: "Old"}, nil)
		mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{}, nil)
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any())
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID, Name: newName}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.UpdateProduct(ctx, productID, req)
		assert.NoError(t, err)
		assert.Equal(t, newName, resp.Name)
	})

	t.Run("UpdateWithStockAndCategories", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		product := products_repo.GetProductWithOptionsRow{ID: productID, Name: "Old", Stock: 10}
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
		mockRepo.EXPECT().CheckCategoryExists(gomock.Any(), int32(1)).Return(true, nil)
		mockRepo.EXPECT().CheckCategoryExists(gomock.Any(), int32(2)).Return(true, nil)
		mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{ID: productID, Name: "New", Stock: 20}, nil)
		mockRepo.EXPECT().ClearProductCategories(gomock.Any(), productID).Return(nil)
		mockRepo.EXPECT().AssignProductCategory(gomock.Any(), gomock.Any()).Return(nil).Times(2)
		mockRepo.EXPECT().CreateStockHistory(gomock.Any(), gomock.Any()).Return(products_repo.StockHistory{}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID, Name: "New", Stock: 20}, nil)

		req := products.UpdateProductRequest{
			Name:        utils.StringPtr("New"),
			Stock:       utils.Int32Ptr(20),
			CategoryIDs: &[]int32{1, 2},
		}
		resp, err := service.UpdateProduct(ctx, productID, req)
		assert.NoError(t, err)
		assert.Equal(t, int32(20), resp.Stock)
	})
}

func TestPrdService_CreateProduct(t *testing.T) {
	ctx := context.Background()
	req := products.CreateProductRequest{
		Name:        "New Product",
		Price:       15000,
		CostPrice:   10000,
		Stock:       50,
		CategoryIDs: []int32{1},
		Options: []products.CreateProductOptionRequest{
			{Name: "Size L", AdditionalPrice: 2000},
		},
	}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockStore := mocks.NewMockStore(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(mockStore, mockRepo, mockLogger, nil, mockActivity)

		productID := uuid.New()
		mockTx, _ := pgxmock.NewConn()

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(pgx.Tx) error) error {
			mockTx.ExpectQuery("INSERT INTO products").
				WithArgs(req.Name, pgxmock.AnyArg(), int64(req.Price), req.Stock, pgxmock.AnyArg()).
				WillReturnRows(pgxmock.NewRows([]string{"id", "name", "image_url", "price", "stock", "created_at", "updated_at", "deleted_at", "cost_price"}).
					AddRow(productID, req.Name, nil, int64(req.Price), req.Stock, time.Now(), time.Now(), nil, pgtype.Numeric{Valid: true}))

			mockTx.ExpectExec("INSERT INTO product_categories").WithArgs(productID, req.CategoryIDs[0]).WillReturnResult(pgxmock.NewResult("INSERT", 1))

			mockTx.ExpectQuery("INSERT INTO product_options").
				WithArgs(productID, req.Options[0].Name, int64(req.Options[0].AdditionalPrice), pgxmock.AnyArg()).
				WillReturnRows(pgxmock.NewRows([]string{"id", "product_id", "name", "additional_price", "image_url", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), productID, req.Options[0].Name, int64(req.Options[0].AdditionalPrice), nil, time.Now(), time.Now(), nil))

			return fn(mockTx)
		})

		mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(products_repo.GetProductByIDRow{ID: productID, Name: req.Name}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.CreateProduct(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NoError(t, mockTx.ExpectationsWereMet())
	})
}

func TestPrdService_ListProducts(t *testing.T) {
	ctx := context.Background()
	req := products.ListProductsRequest{}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return([]products_repo.ListProductsRow{
			{ID: uuid.New(), Name: "P1"},
		}, nil)
		mockRepo.EXPECT().CountProducts(gomock.Any(), gomock.Any()).Return(int64(1), nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.ListProducts(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Products)
	})
}

func TestPrdService_ListDeletedProducts(t *testing.T) {
	ctx := context.Background()
	req := products.ListProductsRequest{}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().ListDeletedProducts(gomock.Any(), gomock.Any()).Return([]products_repo.ListDeletedProductsRow{
			{ID: uuid.New(), Name: "P1 Deleted"},
		}, nil)
		mockRepo.EXPECT().CountDeletedProducts(gomock.Any(), gomock.Any()).Return(int64(1), nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.ListDeletedProducts(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Products)
	})
}

func TestPrdService_DeleteProduct(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{ID: productID, Name: "Del"}, nil)
		mockRepo.EXPECT().SoftDeleteProduct(gomock.Any(), productID).Return(nil)
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any())
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		err := service.DeleteProduct(ctx, productID)
		assert.NoError(t, err)
	})
}

func TestPrdService_RestoreProduct(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{ID: productID}, nil)
		mockRepo.EXPECT().RestoreProduct(gomock.Any(), productID).Return(nil)
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any())
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		err := service.RestoreProduct(ctx, productID)
		assert.NoError(t, err)
	})
}

func TestPrdService_GetStockHistory(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	req := products.ListStockHistoryRequest{}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID}, nil)
		mockRepo.EXPECT().GetStockHistoryByProductWithPagination(gomock.Any(), gomock.Any()).Return([]products_repo.StockHistory{
			{ID: uuid.New(), ProductID: productID, ChangeAmount: 10, CurrentStock: 20},
		}, nil)
		mockRepo.EXPECT().CountStockHistoryByProduct(gomock.Any(), productID).Return(int64(1), nil)

		resp, err := service.GetStockHistory(ctx, productID, req)
		assert.NoError(t, err)
		assert.Len(t, resp.History, 1)
	})
}

func TestPrdService_RestoreProductsBulk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductQuerier(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	mockImgRepo := mocks.NewMockIProductImageRepository(ctrl)
	service := products.NewPrdService(nil, mockRepo, mockLogger, mockImgRepo, mockActivity)

	ctx := context.Background()
	productIDs := []uuid.UUID{uuid.New(), uuid.New()}
	idStrings := []string{productIDs[0].String(), productIDs[1].String()}
	req := products.RestoreBulkRequest{ProductIDs: idStrings}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().RestoreProductsBulk(gomock.Any(), gomock.Any()).Return(nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		err := service.RestoreProductsBulk(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().RestoreProductsBulk(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		err := service.RestoreProductsBulk(ctx, req)
		assert.Error(t, err)
	})

	t.Run("RestoreBulk_RepoError", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().RestoreProductsBulk(gomock.Any(), gomock.Any()).Return(errors.New("repo err"))
		
		err := service.RestoreProductsBulk(ctx, req)
		assert.Error(t, err)
	})

}

func TestPrdService_OptionsCRUD(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	optionID := uuid.New()

	t.Run("CreateOption_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{ID: productID}, nil)
		mockRepo.EXPECT().CreateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID, ProductID: productID}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		req := products.CreateProductOptionRequestStandalone{Name: "Opt"}
		resp, err := service.CreateProductOption(ctx, productID, req)
		assert.NoError(t, err)
		assert.Equal(t, optionID, resp.ID)
	})

	t.Run("UpdateOption_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID, ProductID: productID}, nil)
		mockRepo.EXPECT().UpdateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		req := products.UpdateProductOptionRequest{Name: utils.StringPtr("New Opt")}
		resp, err := service.UpdateProductOption(ctx, productID, optionID, req)
		assert.NoError(t, err)
		assert.Equal(t, optionID, resp.ID)
	})

	t.Run("DeleteOption_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID, ProductID: productID}, nil)
		mockRepo.EXPECT().SoftDeleteProductOption(gomock.Any(), optionID).Return(nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		err := service.DeleteProductOption(ctx, productID, optionID)
		assert.NoError(t, err)
	})
}

func TestPrdService_ListFilters(t *testing.T) {
	ctx := context.Background()
	t.Run("Filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return([]products_repo.ListProductsRow{}, nil)
		mockRepo.EXPECT().CountProducts(gomock.Any(), gomock.Any()).Return(int64(0), nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

		req := products.ListProductsRequest{CategoryID: utils.Int32Ptr(1)}
		_, err := service.ListProducts(ctx, req)
		assert.NoError(t, err)
	})
}

func TestPrdService_StockEdgeCases(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()

	t.Run("GetStockHistory_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.GetStockHistory(ctx, productID, products.ListStockHistoryRequest{})
		assert.ErrorIs(t, err, common.ErrNotFound)
	})
}

func TestPrdService_Lifecycle(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()

	t.Run("RestoreProduct_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{ID: productID}, nil)
		mockRepo.EXPECT().RestoreProduct(gomock.Any(), productID).Return(nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		err := service.RestoreProduct(ctx, productID)
		assert.NoError(t, err)
	})

	t.Run("GetDeletedProduct_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		service := products.NewPrdService(nil, mockRepo, nil, nil, nil)

		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{ID: productID}, nil)

		resp, err := service.GetDeletedProduct(ctx, productID)
		assert.NoError(t, err)
		assert.Equal(t, productID, resp.ID)
	})
}

func TestPrdService_ImageUpload(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	optionID := uuid.New()
	data := []byte("image data")

	t.Run("UploadProductImage_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockImageRepo := mocks.NewMockIProductImageRepository(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImageRepo, mockActivity)

		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID}, nil).Times(2)
		mockImageRepo.EXPECT().UploadImage(gomock.Any(), gomock.Any(), data).Return("http://link", nil)
		mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{ID: productID}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.UploadProductImage(ctx, productID, data)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("UploadProductOptionImage_Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockImageRepo := mocks.NewMockIProductImageRepository(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImageRepo, mockActivity)

		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID, ProductID: productID}, nil)
		mockImageRepo.EXPECT().UploadImage(gomock.Any(), gomock.Any(), data).Return("http://link", nil)
		mockRepo.EXPECT().UpdateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID, ImageUrl: utils.StringPtr("link")}, nil)
		mockImageRepo.EXPECT().PrdImageLink(gomock.Any(), gomock.Any(), gomock.Any()).Return("http://public", nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.UploadProductOptionImage(ctx, productID, optionID, data)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestPrdService_ServiceErrorCases(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	optionID := uuid.New()

	t.Run("GetProductByID_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.GetProductByID(ctx, productID)
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("UpdateProductOption_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.UpdateProductOption(ctx, productID, optionID, products.UpdateProductOptionRequest{})
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("DeleteProductOption_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		err := service.DeleteProductOption(ctx, productID, optionID)
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("GetDeletedProduct_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.GetDeletedProduct(ctx, productID)
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("DeleteProduct_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		err := service.DeleteProduct(ctx, productID)
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("RestoreProduct_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		err := service.RestoreProduct(ctx, productID)
		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("UpdateProduct_StockHistoryDetails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		product := products_repo.GetProductWithOptionsRow{ID: productID, Name: "Old", Stock: 10}
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
		mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{ID: productID, Name: "Old", Stock: 5}, nil)
		// changeAmount = -5, so changeType will be Correction (default)
		mockRepo.EXPECT().CreateStockHistory(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx any, params products_repo.CreateStockHistoryParams) (products_repo.StockHistory, error) {
			assert.Equal(t, int32(-5), params.ChangeAmount)
			assert.Equal(t, products_repo.StockChangeTypeCorrection, params.ChangeType)
			assert.Equal(t, "Correction note", *params.Note)
			return products_repo.StockHistory{}, nil
		})
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID, Name: "Old", Stock: 5}, nil)

		req := products.UpdateProductRequest{
			Stock: utils.Int32Ptr(5),
			Note:  utils.StringPtr("Correction note"),
		}
		_, err := service.UpdateProduct(ctx, productID, req)
		assert.NoError(t, err)
	})

	t.Run("UpdateProduct_StockHistoryRestock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		product := products_repo.GetProductWithOptionsRow{ID: productID, Name: "Old", Stock: 10}
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
		mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{ID: productID, Name: "Old", Stock: 15}, nil)
		// changeAmount = 5, so changeType will be Restock (if not provided)
		mockRepo.EXPECT().CreateStockHistory(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx any, params products_repo.CreateStockHistoryParams) (products_repo.StockHistory, error) {
			assert.Equal(t, int32(5), params.ChangeAmount)
			assert.Equal(t, products_repo.StockChangeTypeRestock, params.ChangeType)
			return products_repo.StockHistory{}, nil
		})
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID, Name: "Old", Stock: 15}, nil)

		req := products.UpdateProductRequest{
			Stock: utils.Int32Ptr(15),
		}
		_, err := service.UpdateProduct(ctx, productID, req)
		assert.NoError(t, err)
	})

	t.Run("UpdateProduct_StockHistoryManualChangeType", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, mockActivity)

		product := products_repo.GetProductWithOptionsRow{ID: productID, Name: "Old", Stock: 10}
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
		mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{ID: productID, Name: "Old", Stock: 20}, nil)
		mockRepo.EXPECT().CreateStockHistory(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx any, params products_repo.CreateStockHistoryParams) (products_repo.StockHistory, error) {
			assert.Equal(t, products_repo.StockChangeTypeCorrection, params.ChangeType)
			return products_repo.StockHistory{}, nil
		})
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), productID.String(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID, Name: "Old", Stock: 20}, nil)

		req := products.UpdateProductRequest{
			Stock:      utils.Int32Ptr(20),
			ChangeType: utils.StringPtr(string(products_repo.StockChangeTypeCorrection)),
		}
		_, err := service.UpdateProduct(ctx, productID, req)
		assert.NoError(t, err)
	})

	t.Run("GetProductByID_ImageError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockImageRepo := mocks.NewMockIProductImageRepository(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImageRepo, nil)

		imageUrl := "products/test.jpg"
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{
			ID:       productID,
			Name:     "Test Product",
			ImageUrl: &imageUrl,
		}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockImageRepo.EXPECT().PrdImageLink(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("link fail"))

		_, err := service.GetProductByID(ctx, productID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get main image link")
	})

	t.Run("GetProductByID_OptionImageError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockImageRepo := mocks.NewMockIProductImageRepository(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImageRepo, nil)

		optionID := uuid.New()
		optImg := "options/opt.jpg"
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{
			ID:   productID,
			Name: "Test Product",
			Options: []map[string]interface{}{{
				"id":        optionID,
				"name":      "Opt",
				"image_url": optImg,
			}},
		}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockImageRepo.EXPECT().PrdImageLink(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("link fail")).AnyTimes()

		_, err := service.GetProductByID(ctx, productID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get option image link")
	})

	t.Run("DeleteProductOption_DBError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		err := service.DeleteProductOption(ctx, productID, optionID)
		assert.Error(t, err)
	})

	t.Run("DeleteProductOption_SoftDeleteError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{Name: "Opt"}, nil)
		mockRepo.EXPECT().SoftDeleteProductOption(gomock.Any(), optionID).Return(errors.New("soft delete error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		err := service.DeleteProductOption(ctx, productID, optionID)
		assert.Error(t, err)
	})

	t.Run("UpdateProductOption_DBError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.UpdateProductOption(ctx, productID, optionID, products.UpdateProductOptionRequest{})
		assert.Error(t, err)
	})

	t.Run("UpdateProductOption_UpdateError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID}, nil)
		mockRepo.EXPECT().UpdateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, errors.New("update error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.UpdateProductOption(ctx, productID, optionID, products.UpdateProductOptionRequest{})
		assert.Error(t, err)
	})

	t.Run("UploadProductOptionImage_FileTooLarge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID}, nil)
		largeData := make([]byte, 3*1024*1024) // 3MB > 2MB limit
		_, err := service.UploadProductOptionImage(ctx, productID, optionID, largeData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds the limit")
	})

	t.Run("UploadProductOptionImage_UploadError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockImageRepo := mocks.NewMockIProductImageRepository(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImageRepo, nil)
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{ID: optionID}, nil)
		mockImageRepo.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("upload fail"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.UploadProductOptionImage(ctx, productID, optionID, []byte("data"))
		assert.Error(t, err)
	})

	t.Run("UpdateProduct_CategoryNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{ID: productID}, nil)
		mockRepo.EXPECT().CheckCategoryExists(gomock.Any(), int32(99)).Return(false, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.UpdateProduct(ctx, productID, products.UpdateProductRequest{CategoryIDs: &[]int32{99}})
		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
	})

	t.Run("UpdateProduct_CheckCategoryError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{ID: productID}, nil)
		mockRepo.EXPECT().CheckCategoryExists(gomock.Any(), int32(1)).Return(false, errors.New("check error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		_, err := service.UpdateProduct(ctx, productID, products.UpdateProductRequest{CategoryIDs: &[]int32{1}})
		assert.Error(t, err)
	})

	t.Run("UploadProductImage_FileTooLarge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		mockRepo.EXPECT().GetProductByID(gomock.Any(), productID).Return(products_repo.GetProductByIDRow{ID: productID}, nil)
		largeData := make([]byte, 6*1024*1024) // 6MB > 5MB limit
		_, err := service.UploadProductImage(ctx, productID, largeData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds the limit")
	})

	t.Run("GetDeletedProduct_DBError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		_, err := service.GetDeletedProduct(ctx, productID)
		assert.Error(t, err)
	})

	t.Run("GetDeletedProduct_WithOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		service := products.NewPrdService(nil, mockRepo, nil, nil, nil)

		optionID := uuid.New()
		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{
			ID:   productID,
			Name: "Deleted P",
			Options: []map[string]interface{}{{
				"id":   optionID,
				"name": "Opt 1",
			}},
			DeletedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}, nil)

		resp, err := service.GetDeletedProduct(ctx, productID)
		assert.NoError(t, err)
		assert.Equal(t, "Deleted P", resp.Name)
		assert.Len(t, resp.Options, 1)
		assert.NotNil(t, resp.DeletedAt)
	})

	t.Run("RestoreProduct_RepoError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().GetDeletedProduct(gomock.Any(), productID).Return(products_repo.GetDeletedProductRow{ID: productID}, nil)
		mockRepo.EXPECT().RestoreProduct(gomock.Any(), productID).Return(errors.New("restore fail"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		err := service.RestoreProduct(ctx, productID)
		assert.Error(t, err)
	})

	t.Run("CreateProduct_TxError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockStore := mocks.NewMockStore(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(mockStore, mockRepo, mockLogger, nil, nil)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(errors.New("tx fail"))

		_, err := service.CreateProduct(ctx, products.CreateProductRequest{Name: "P"})
		assert.Error(t, err)
	})

	t.Run("ListProducts_RepoError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		mockRepo.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return(nil, errors.New("list fail"))
		mockRepo.EXPECT().CountProducts(gomock.Any(), gomock.Any()).Return(int64(0), nil).AnyTimes()
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		_, err := service.ListProducts(ctx, products.ListProductsRequest{})
		assert.Error(t, err)
	})

	t.Run("UpdateProduct_UpdatesErrors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)

		product := products_repo.GetProductWithOptionsRow{ID: productID, Name: "Old"}
		
		t.Run("ClearCategoriesFail", func(t *testing.T) {
			mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
			mockRepo.EXPECT().CheckCategoryExists(gomock.Any(), gomock.Any()).Return(true, nil)
			mockRepo.EXPECT().ClearProductCategories(gomock.Any(), productID).Return(errors.New("clear fail"))
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			_, err := service.UpdateProduct(ctx, productID, products.UpdateProductRequest{CategoryIDs: &[]int32{1}})
			assert.Error(t, err)
		})

		t.Run("AssignCategoryFail", func(t *testing.T) {
			mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
			mockRepo.EXPECT().CheckCategoryExists(gomock.Any(), gomock.Any()).Return(true, nil)
			mockRepo.EXPECT().ClearProductCategories(gomock.Any(), productID).Return(nil)
			mockRepo.EXPECT().AssignProductCategory(gomock.Any(), gomock.Any()).Return(errors.New("assign fail"))
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			_, err := service.UpdateProduct(ctx, productID, products.UpdateProductRequest{CategoryIDs: &[]int32{1}})
			assert.Error(t, err)
		})

		t.Run("RepoUpdateFail", func(t *testing.T) {
			mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(product, nil)
			mockRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).Return(products_repo.Product{}, errors.New("update fail"))
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			_, err := service.UpdateProduct(ctx, productID, products.UpdateProductRequest{Name: utils.StringPtr("New")})
			assert.Error(t, err)
		})
	})

	t.Run("CreateOption_ParentNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		productID := uuid.New()
		req := products.CreateProductOptionRequestStandalone{Name: "Option 1"}
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		
		res, err := service.CreateProductOption(ctx, productID, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("CreateOption_RepoError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		productID := uuid.New()
		req := products.CreateProductOptionRequestStandalone{Name: "Option 1"}
		mockRepo.EXPECT().GetProductWithOptions(gomock.Any(), productID).Return(products_repo.GetProductWithOptionsRow{}, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().CreateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, errors.New("repo err"))
		
		res, err := service.CreateProductOption(ctx, productID, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("UpdateOption_RepoError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, nil, nil)
		productID := uuid.New()
		optionID := uuid.New()
		req := products.UpdateProductOptionRequest{Name: utils.StringPtr("New Name")}
		mockRepo.EXPECT().GetProductOption(gomock.Any(), products_repo.GetProductOptionParams{ID: optionID, ProductID: productID}).Return(products_repo.ProductOption{}, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		mockRepo.EXPECT().UpdateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, errors.New("repo err"))
		
		res, err := service.UpdateProductOption(ctx, productID, optionID, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("UploadOptionImage_UploadError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockImgRepo := mocks.NewMockIProductImageRepository(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImgRepo, nil)
		productID := uuid.New()
		optionID := uuid.New()
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, nil)
		mockImgRepo.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("upload fail"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		
		res, err := service.UploadProductOptionImage(ctx, productID, optionID, []byte("data"))
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("UploadOptionImage_RepoError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockProductQuerier(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockImgRepo := mocks.NewMockIProductImageRepository(ctrl)
		service := products.NewPrdService(nil, mockRepo, mockLogger, mockImgRepo, nil)
		productID := uuid.New()
		optionID := uuid.New()
		mockRepo.EXPECT().GetProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, nil)
		mockImgRepo.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return("url", nil)
		mockRepo.EXPECT().UpdateProductOption(gomock.Any(), gomock.Any()).Return(products_repo.ProductOption{}, errors.New("repo err"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		
		res, err := service.UploadProductOptionImage(ctx, productID, optionID, []byte("data"))
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
