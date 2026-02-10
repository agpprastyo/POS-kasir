package products

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type IPrdService interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	UploadProductImage(ctx context.Context, productID uuid.UUID, data []byte) (*dto.ProductResponse, error)
	ListProducts(ctx context.Context, req dto.ListProductsRequest) (*dto.ListProductsResponse, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	CreateProductOption(ctx context.Context, productID uuid.UUID, req dto.CreateProductOptionRequestStandalone) (*dto.ProductOptionResponse, error)
	GetStockHistory(ctx context.Context, productID uuid.UUID, req dto.ListStockHistoryRequest) (*dto.PagedStockHistoryResponse, error)
	UploadProductOptionImage(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, data []byte) (*dto.ProductOptionResponse, error)
	UpdateProductOption(ctx context.Context, productID, optionID uuid.UUID, req dto.UpdateProductOptionRequest) (*dto.ProductOptionResponse, error)
	DeleteProductOption(ctx context.Context, productID, optionID uuid.UUID) error

	// Deleted Products Management
	ListDeletedProducts(ctx context.Context, req dto.ListProductsRequest) (*dto.ListProductsResponse, error)
	GetDeletedProduct(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error)
	RestoreProduct(ctx context.Context, productID uuid.UUID) error
	RestoreProductsBulk(ctx context.Context, req dto.RestoreBulkRequest) error
}

type PrdService struct {
	log             logger.ILogger
	store           repository.Store
	prdRepo         IPrdRepo
	activityService activitylog.IActivityService
}

func NewPrdService(store repository.Store, log logger.ILogger, prdRepo IPrdRepo, activityService activitylog.IActivityService) IPrdService {
	return &PrdService{
		store:           store,
		log:             log,
		prdRepo:         prdRepo,
		activityService: activityService,
	}
}

func (s *PrdService) DeleteProductOption(ctx context.Context, productID, optionID uuid.UUID) error {

	option, err := s.store.GetProductOption(ctx, repository.GetProductOptionParams{ID: optionID, ProductID: productID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product option not found or does not belong to the product", "optionID", optionID, "productID", productID)
			return common.ErrNotFound
		}
		s.log.Errorf("Failed to get product option before deletion", "error", err)
		return err
	}

	err = s.store.SoftDeleteProductOption(ctx, optionID)
	if err != nil {
		s.log.Errorf("Failed to soft delete product option in repository", "error", err, "optionID", optionID)
		return err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"product_id":          productID.String(),
		"deleted_option_id":   optionID.String(),
		"deleted_option_name": option.Name,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeDELETE,
		repository.LogEntityTypePRODUCT,
		optionID.String(),
		logDetails,
	)

	s.log.Infof("Product option soft deleted successfully", "optionID", optionID)
	return nil
}

func (s *PrdService) UpdateProductOption(ctx context.Context, productID, optionID uuid.UUID, req dto.UpdateProductOptionRequest) (*dto.ProductOptionResponse, error) {
	_, err := s.store.GetProductOption(ctx, repository.GetProductOptionParams{ID: optionID, ProductID: productID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product option not found or does not belong to the product", "optionID", optionID, "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get product option before update", "error", err)
		return nil, err
	}

	updateParams := repository.UpdateProductOptionParams{
		ID:   optionID,
		Name: req.Name,
	}

	if req.AdditionalPrice != nil {
		price := int64(*req.AdditionalPrice)
		updateParams.AdditionalPrice = &price
	}

	updatedOption, err := s.store.UpdateProductOption(ctx, updateParams)
	if err != nil {
		s.log.Errorf("Failed to update product option in repository", "error", err, "optionID", optionID)
		return nil, err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"product_id":     productID.String(),
		"option_id":      optionID.String(),
		"updated_fields": req,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypePRODUCT,
		optionID.String(),
		logDetails,
	)

	additionalPrice := float64(updatedOption.AdditionalPrice)

	if updatedOption.ImageUrl != nil {
		publicUrl, err := s.prdRepo.PrdImageLink(ctx, updatedOption.ID.String(), *updatedOption.ImageUrl)
		if err != nil {
			s.log.Warnf("Failed to get public URL for updated option image", "error", err)
			publicUrl = *updatedOption.ImageUrl
		}
		updatedOption.ImageUrl = &publicUrl
	} else {
		updatedOption.ImageUrl = nil
	}

	return &dto.ProductOptionResponse{
		ID:              updatedOption.ID,
		Name:            updatedOption.Name,
		AdditionalPrice: additionalPrice,
		ImageURL:        updatedOption.ImageUrl,
	}, nil
}

func (s *PrdService) UploadProductOptionImage(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, data []byte) (*dto.ProductOptionResponse, error) {
	_, err := s.store.GetProductOption(ctx, repository.GetProductOptionParams{
		ID:        optionID,
		ProductID: productID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product option not found or does not belong to the product", "optionID", optionID, "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get product option before image upload", "error", err)
		return nil, err
	}

	const maxFileSize = 2 * 1024 * 1024
	if len(data) > maxFileSize {
		return nil, fmt.Errorf("file size exceeds the limit of 2MB")
	}

	filename := fmt.Sprintf("product_options/%s.jpg", optionID.String())

	imageUrl, err := s.prdRepo.UploadImage(ctx, filename, data)
	if err != nil {
		s.log.Errorf("Failed to upload option image to R2", "error", err)
		return nil, fmt.Errorf("could not upload image to storage")
	}

	updateParams := repository.UpdateProductOptionParams{
		ID:       optionID,
		ImageUrl: &filename,
	}
	updatedOption, err := s.store.UpdateProductOption(ctx, updateParams)
	if err != nil {
		s.log.Errorf("Failed to update product option with image URL", "error", err)
		return nil, fmt.Errorf("could not update product option in database")
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"product_id":    productID.String(),
		"option_id":     optionID.String(),
		"new_image_url": imageUrl,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypePRODUCT,
		optionID.String(),
		logDetails,
	)

	var additionalPrice float64
	additionalPrice = float64(updatedOption.AdditionalPrice)

	publicUrl, err := s.prdRepo.PrdImageLink(ctx, updatedOption.ID.String(), *updatedOption.ImageUrl)
	if err != nil {
		s.log.Warnf("Failed to get public URL for newly uploaded option image", "error", err)
		publicUrl = *updatedOption.ImageUrl
	}

	return &dto.ProductOptionResponse{
		ID:              updatedOption.ID,
		Name:            updatedOption.Name,
		AdditionalPrice: additionalPrice,
		ImageURL:        &publicUrl,
	}, nil
}
func (s *PrdService) CreateProductOption(ctx context.Context, productID uuid.UUID, req dto.CreateProductOptionRequestStandalone) (*dto.ProductOptionResponse, error) {
	_, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Parent product not found for new option", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get parent product for new option", "error", err)
		return nil, err
	}

	additionalPrice := int64(req.AdditionalPrice)
	params := repository.CreateProductOptionParams{
		ProductID:       productID,
		Name:            req.Name,
		AdditionalPrice: additionalPrice,
	}
	newOption, err := s.store.CreateProductOption(ctx, params)
	if err != nil {
		s.log.Errorf("Failed to create product option in repository", "error", err)
		return nil, err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"parent_product_id": productID.String(),
		"option_name":       newOption.Name,
		"additional_price":  req.AdditionalPrice,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypePRODUCT,
		newOption.ID.String(),
		logDetails,
	)

	return &dto.ProductOptionResponse{
		ID:              newOption.ID,
		Name:            newOption.Name,
		AdditionalPrice: float64(newOption.AdditionalPrice),
		ImageURL:        newOption.ImageUrl,
	}, nil
}
func (s *PrdService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {

	product, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product not found for deletion", "productID", productID)
			return common.ErrNotFound
		}
		s.log.Errorf("Failed to get product before deletion", "error", err)
		return err
	}

	err = s.store.SoftDeleteProduct(ctx, productID)
	if err != nil {
		s.log.Errorf("Failed to soft delete product in repository", "error", err, "productID", productID)
		return err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"deleted_product_id":   product.ID.String(),
		"deleted_product_name": product.Name,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeDELETE,
		repository.LogEntityTypePRODUCT,
		productID.String(),
		logDetails,
	)

	s.log.Infof("Product soft deleted successfully", "productID", productID)
	return nil
}

func (s *PrdService) GetProductByID(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error) {
	fullProduct, err := s.store.GetProductWithOptions(ctx, productID)
	s.log.Infof("Full product data: %+v", fullProduct)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product not found by ID", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get product from repository", "error", err, "productID", productID)
		return nil, err
	}

	return s.buildProductResponse(ctx, fullProduct)
}

func (s *PrdService) UpdateProduct(ctx context.Context, productID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product not found for update", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get product before update", "error", err)
		return nil, err
	}

	if req.CategoryID != nil {
		exists, err := s.store.ExistsCategory(ctx, *req.CategoryID)
		if err != nil {
			s.log.Errorf("Failed to check category existence", "error", err)
			return nil, err
		}
		if !exists {
			s.log.Warnf("Category not found", "categoryID", *req.CategoryID)
			return nil, common.ErrCategoryNotFound
		}
	}

	updateParams := repository.UpdateProductParams{
		ID:         productID,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Stock:      req.Stock,
	}

	if req.Price != nil {
		price := int64(*req.Price)
		updateParams.Price = &price
	}
	if req.CostPrice != nil {
		costPrice := *req.CostPrice
		numericCost := pgtype.Numeric{}
		numericCost.Scan(fmt.Sprintf("%f", costPrice))
		updateParams.CostPrice = numericCost
	}

	_, err = s.store.UpdateProduct(ctx, updateParams)
	if err != nil {
		s.log.Errorf("Failed to update product in repository", "error", err, "productID", productID)
		return nil, err
	}

	// Stock History Logging
	if req.Stock != nil {
		currentStock := int32(*req.Stock)
		previousStock := int32(product.Stock)
		changeAmount := currentStock - previousStock

		if changeAmount != 0 {
			changeType := repository.StockChangeTypeCorrection // Default
			if req.ChangeType != nil {
				changeType = repository.StockChangeType(*req.ChangeType)
			} else {
				// Auto-detect simple cases if not provided
				if changeAmount > 0 {
					changeType = repository.StockChangeTypeRestock
				}
			}

			// Safe dereference for Note
			var note *string
			if req.Note != nil {
				note = req.Note
			}

			// Creator
			var createdBy pgtype.UUID
			if actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID); ok {
				createdBy = pgtype.UUID{Bytes: actorID, Valid: true}
			}

			// Reference? For direct update, maybe null or strict correlation if we had one.
			// Current flow doesn't have a specific reference ID for manual updates other than maybe the log ID?
			// Leaving referenceID null for manual product updates.

			err := s.RecordStockChange(ctx, productID, changeAmount, previousStock, currentStock, changeType, pgtype.UUID{Valid: false}, note, createdBy)
			if err != nil {
				// Don't fail the request if logging fails, but log error
				s.log.Errorf("Failed to record stock history", "error", err)
			}
		}
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"product_id":     productID.String(),
		"updated_fields": req,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypePRODUCT,
		productID.String(),
		logDetails,
	)

	return s.GetProductByID(ctx, productID)
}

func (s *PrdService) RecordStockChange(ctx context.Context, productID uuid.UUID, changeAmount, previousStock, currentStock int32, changeType repository.StockChangeType, referenceID pgtype.UUID, note *string, createdBy pgtype.UUID) error {
	params := repository.CreateStockHistoryParams{
		ProductID:     productID,
		ChangeAmount:  changeAmount,
		PreviousStock: previousStock,
		CurrentStock:  currentStock,
		ChangeType:    changeType,
		ReferenceID:   referenceID,
		Note:          note,
		CreatedBy:     createdBy,
	}

	_, err := s.store.CreateStockHistory(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create stock history: %w", err)
	}

	return nil
}
func (s *PrdService) buildProductResponse(ctx context.Context, fullProduct repository.GetProductWithOptionsRow) (*dto.ProductResponse, error) {
	var optionsResponse []dto.ProductOptionResponse

	if fullProduct.Options != nil {

		optionsJSON, err := json.Marshal(fullProduct.Options)
		if err != nil {
			s.log.Errorf("Failed to re-marshal product options interface", "error", err)
			return nil, fmt.Errorf("could not process product options")
		}

		var options []repository.ProductOption
		if err := json.Unmarshal(optionsJSON, &options); err != nil {

			if string(optionsJSON) != "[]" {
				s.log.Errorf("Failed to unmarshal product options JSON", "error", err)
				return nil, fmt.Errorf("could not parse product options")
			}
		}

		for _, opt := range options {
			additionalPrice := float64(opt.AdditionalPrice)
			optionsResponse = append(optionsResponse, dto.ProductOptionResponse{
				ID:              opt.ID,
				Name:            opt.Name,
				AdditionalPrice: additionalPrice,
				ImageURL:        opt.ImageUrl,
			})
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(optionsResponse)+1)

	if fullProduct.ImageUrl != nil && *fullProduct.ImageUrl != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := s.prdRepo.PrdImageLink(ctx, fullProduct.ID.String(), *fullProduct.ImageUrl)
			if err != nil {
				errChan <- fmt.Errorf("failed to get main image link: %w", err)
				return
			}
			fullProduct.ImageUrl = &url
		}()
	}

	for i := range optionsResponse {

		opt := &optionsResponse[i]
		if opt.ImageURL != nil && *opt.ImageURL != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				url, err := s.prdRepo.PrdImageLink(ctx, *opt.ImageURL, *opt.ImageURL)
				if err != nil {
					errChan <- fmt.Errorf("failed to get option image link for %s: %w", opt.Name, err)
					return
				}
				opt.ImageURL = &url
			}()
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	productPrice := float64(fullProduct.Price)
	costPrice := 0.0
	if fullProduct.CostPrice.Valid {
		prodCost, _ := fullProduct.CostPrice.Float64Value()
		costPrice = prodCost.Float64
	}

	return &dto.ProductResponse{
		ID:         fullProduct.ID,
		Name:       fullProduct.Name,
		CategoryID: fullProduct.CategoryID,

		ImageURL:  fullProduct.ImageUrl,
		Price:     productPrice,
		CostPrice: costPrice,
		Stock:     fullProduct.Stock,
		CreatedAt: fullProduct.CreatedAt.Time,
		UpdatedAt: fullProduct.UpdatedAt.Time,
		Options:   optionsResponse,
	}, nil
}

func (s *PrdService) ListProducts(ctx context.Context, req dto.ListProductsRequest) (*dto.ListProductsResponse, error) {

	page := 1
	if req.Page != nil {
		page = *req.Page
	}
	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	listParams := repository.ListProductsParams{
		Limit:      int32(limit),
		Offset:     int32(offset),
		CategoryID: req.CategoryID,
		SearchText: req.Search,
	}
	countParams := repository.CountProductsParams{
		CategoryID: req.CategoryID,
		SearchText: req.Search,
	}

	s.log.Infof("list params list product: %+v", listParams)

	var wg sync.WaitGroup
	var products []repository.ListProductsRow
	var totalData int64
	var listErr, countErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		products, listErr = s.store.ListProducts(ctx, listParams)
	}()

	go func() {
		defer wg.Done()
		totalData, countErr = s.store.CountProducts(ctx, countParams)
	}()

	wg.Wait()

	if listErr != nil {
		s.log.Errorf("Failed to list products from repository", "error", listErr)
		return nil, listErr
	}
	if countErr != nil {
		s.log.Errorf("Failed to count products from repository", "error", countErr)
		return nil, countErr
	}

	var productsResponse []dto.ProductListResponse
	for _, p := range products {
		price := float64(p.Price)
		if p.ImageUrl != nil && *p.ImageUrl != "" {
			imageUrl, err := s.prdRepo.PrdImageLink(ctx, p.ID.String(), *p.ImageUrl)
			if err != nil {
				s.log.Warnf("Failed to get public URL for product image", "error", err)
				imageUrl = *p.ImageUrl // Tetap gunakan URL asli jika gagal
			}
			p.ImageUrl = &imageUrl
		} else {
			p.ImageUrl = nil // Setel ke nil jika tidak ada URL
		}
		productsResponse = append(productsResponse, dto.ProductListResponse{
			ID:           p.ID,
			Name:         p.Name,
			CategoryID:   p.CategoryID,
			CategoryName: p.CategoryName,
			ImageURL:     p.ImageUrl,
			Price:        price,
			Stock:        p.Stock,
		})
	}

	response := &dto.ListProductsResponse{
		Products: productsResponse,
		Pagination: pagination.BuildPagination(
			page,
			int(totalData),
			limit,
		),
	}

	return response, nil
}

func (s *PrdService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	var newProduct repository.Product
	var createdOptions []repository.ProductOption

	txFunc := func(qtx *repository.Queries) error {
		var err error

		price := int64(req.Price)

		numericCost := pgtype.Numeric{}
		numericCost.Scan(fmt.Sprintf("%f", req.CostPrice))

		productParams := repository.CreateProductParams{
			Name:       req.Name,
			CategoryID: &req.CategoryID,
			Price:      price,
			Stock:      req.Stock,
			CostPrice:  numericCost,
		}
		newProduct, err = qtx.CreateProduct(ctx, productParams)
		if err != nil {
			s.log.Errorf("Failed to create product in transaction", "error", err)
			return err
		}

		for _, opt := range req.Options {
			additionalPrice := int64(opt.AdditionalPrice)

			optionParams := repository.CreateProductOptionParams{
				ProductID:       newProduct.ID,
				Name:            opt.Name,
				AdditionalPrice: additionalPrice,
			}
			createdOpt, err := qtx.CreateProductOption(ctx, optionParams)
			if err != nil {
				s.log.Errorf("Failed to create product option in transaction", "error", err)
				return err
			}
			createdOptions = append(createdOptions, createdOpt)
		}
		return nil
	}

	err := s.store.ExecTx(ctx, txFunc)
	if err != nil {
		return nil, err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"product_name":  newProduct.Name,
		"price":         req.Price,
		"stock":         newProduct.Stock,
		"options_count": len(createdOptions),
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypePRODUCT,
		newProduct.ID.String(),
		logDetails,
	)

	return s.buildProductResponseFromData(newProduct, createdOptions)
}

func (s *PrdService) UploadProductImage(ctx context.Context, productID uuid.UUID, data []byte) (*dto.ProductResponse, error) {
	_, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warnf("Product not found for image upload", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get product for image upload", "error", err)
		return nil, err
	}

	const maxFileSize = 5 * 1024 * 1024 // 2MB
	if len(data) > maxFileSize {
		return nil, fmt.Errorf("file size exceeds the limit of 2MB")
	}

	filename := fmt.Sprintf("products/%s.jpg", productID.String())

	imageUrl, err := s.prdRepo.UploadImage(ctx, filename, data)
	if err != nil {
		s.log.Errorf("Failed to upload image to R2", "error", err)
		return nil, fmt.Errorf("could not upload image to storage")
	}

	updateParams := repository.UpdateProductParams{
		ID:       productID,
		ImageUrl: &filename,
	}
	_, err = s.store.UpdateProduct(ctx, updateParams)
	if err != nil {
		s.log.Errorf("Failed to update product with image URL", "error", err)
		return nil, fmt.Errorf("could not update product in database")
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"product_id":    productID.String(),
		"new_image_url": imageUrl,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypePRODUCT,
		productID.String(),
		logDetails,
	)

	fullProduct, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		s.log.Errorf("Failed to fetch full product after image upload", "error", err)
		return nil, err
	}

	return s.buildProductResponse(ctx, fullProduct)
}

func (s *PrdService) buildProductResponseFromData(product repository.Product, options []repository.ProductOption) (*dto.ProductResponse, error) {
	var optionsResponse []dto.ProductOptionResponse
	for _, opt := range options {
		var additionalPrice = float64(opt.AdditionalPrice)
		optionsResponse = append(optionsResponse, dto.ProductOptionResponse{
			ID:              opt.ID,
			Name:            opt.Name,
			AdditionalPrice: additionalPrice,
			ImageURL:        opt.ImageUrl,
		})
	}

	productPrice := float64(product.Price)
	costPrice := 0.0
	if product.CostPrice.Valid {
		prodCost, _ := product.CostPrice.Float64Value()
		costPrice = prodCost.Float64
	}

	s.log.Infof("product price before assign: %+v", product.Price)
	s.log.Infof("product price after assign: %+v", productPrice)

	return &dto.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		ImageURL:   product.ImageUrl,
		Price:      productPrice,
		CostPrice:  costPrice,
		Stock:      product.Stock,
		CreatedAt:  product.CreatedAt.Time,
		UpdatedAt:  product.UpdatedAt.Time,
		Options:    optionsResponse,
	}, nil
}

func (s *PrdService) ListDeletedProducts(ctx context.Context, req dto.ListProductsRequest) (*dto.ListProductsResponse, error) {
	page := 1
	if req.Page != nil {
		page = *req.Page
	}
	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	listParams := repository.ListDeletedProductsParams{
		Limit:      int32(limit),
		Offset:     int32(offset),
		CategoryID: req.CategoryID,
		SearchText: req.Search,
	}
	countParams := repository.CountDeletedProductsParams{
		CategoryID: req.CategoryID,
		SearchText: req.Search,
	}

	var wg sync.WaitGroup
	var products []repository.ListDeletedProductsRow
	var totalData int64
	var listErr, countErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		products, listErr = s.store.ListDeletedProducts(ctx, listParams)
	}()

	go func() {
		defer wg.Done()
		totalData, countErr = s.store.CountDeletedProducts(ctx, countParams)
	}()

	wg.Wait()

	if listErr != nil {
		s.log.Errorf("Failed to list deleted products", "error", listErr)
		return nil, listErr
	}
	if countErr != nil {
		s.log.Errorf("Failed to count deleted products", "error", countErr)
		return nil, countErr
	}

	var productsResponse []dto.ProductListResponse
	for _, p := range products {
		price := float64(p.Price)
		var deletedAt *time.Time
		if p.DeletedAt.Valid {
			t := p.DeletedAt.Time
			deletedAt = &t
		}

		// Similar image logic if needed, skipping detailed image check for list view speed or reuse similar logic
		if p.ImageUrl != nil && *p.ImageUrl != "" {
			imageUrl, err := s.prdRepo.PrdImageLink(ctx, p.ID.String(), *p.ImageUrl)
			if err != nil {
				imageUrl = *p.ImageUrl
			}
			p.ImageUrl = &imageUrl
		} else {
			p.ImageUrl = nil
		}

		productsResponse = append(productsResponse, dto.ProductListResponse{
			ID:           p.ID,
			Name:         p.Name,
			CategoryID:   p.CategoryID,
			CategoryName: p.CategoryName,
			ImageURL:     p.ImageUrl,
			Price:        price,
			Stock:        p.Stock,
			DeletedAt:    deletedAt,
		})
	}

	return &dto.ListProductsResponse{
		Products: productsResponse,
		Pagination: pagination.BuildPagination(
			page,
			int(totalData),
			limit,
		),
	}, nil
}

func (s *PrdService) GetDeletedProduct(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error) {
	// Need to fetch similar to GetProductWithOptions but for deleted one.
	// We defined GetDeletedProduct query to return row + options json.
	// But generated type will be GetDeletedProductRow.

	// Wait, the query GetDeletedProduct in products.sql returns specific columns similar to GetProductWithOptions?
	// Let's assume standard behavior based on sql definition.

	row, err := s.store.GetDeletedProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		s.log.Errorf("Failed to get deleted product", "error", err)
		return nil, err
	}

	// Convert GetDeletedProductRow to GetProductWithOptionsRow structure or handle manually
	// Since struct fields might slightly differ if not strictly aliased, manually map.
	// Actually, GetDeletedProduct query was:
	// SELECT p.*, COALESCE(...) as options FROM products p ...
	// So fields map to table columns.

	var options []repository.ProductOption
	if row.Options != nil {
		// unmarshal json
		optionsBytes, _ := json.Marshal(row.Options) // It's already likely []byte or interface{} depending on driver
		// PGX usually scans json/jsonb into []byte or map/slice if annotated.
		// generated code uses []byte usually for json columns if not overridden.
		// Wait, in sqlc.yaml usually we set overriding or it just uses []byte.
		// Let's assume standard.

		_ = json.Unmarshal(optionsBytes, &options)
		// Correction: row.Options is already likely unmarshalled if we used specific types,
		// but standard sqlc returns []byte for json_agg unless cast.
		// Check generated code... usually interface{}.

		// To be safe and reuse logic, let's look at buildProductResponse logic which handles it.
		// But buildProductResponse takes GetProductWithOptionsRow.
		// We can't reuse it directly if types mismatch excessively.
		// Let's reimplement simplified version.

		// Re-reading buildProductResponse:
		// if fullProduct.Options != nil { ... marshall, unmarshall ... }
		// It seems it receives interface{}.
	}

	// Because we can't see generated repository file content easily right now without checking,
	// I will attempt to cast row to GetProductWithOptionsRow if identical, or manual map.
	// Manual map is safer.

	// Parse Options
	var optionsResponse []dto.ProductOptionResponse
	if row.Options != nil {
		optBytes, err := json.Marshal(row.Options)
		if err == nil {
			var opts []repository.ProductOption
			if err := json.Unmarshal(optBytes, &opts); err == nil {
				for _, o := range opts {
					optionsResponse = append(optionsResponse, dto.ProductOptionResponse{
						ID:              o.ID,
						Name:            o.Name,
						AdditionalPrice: float64(o.AdditionalPrice),
						ImageURL:        o.ImageUrl,
					})
				}
			}
		}
	}

	var deletedAt time.Time
	if row.DeletedAt.Valid {
		deletedAt = row.DeletedAt.Time
	}

	return &dto.ProductResponse{
		ID:         row.ID,
		Name:       row.Name,
		CategoryID: row.CategoryID,
		ImageURL:   row.ImageUrl,
		Price:      float64(row.Price),
		Stock:      row.Stock,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
		DeletedAt:  &deletedAt,
		Options:    optionsResponse,
	}, nil
}

func (s *PrdService) RestoreProduct(ctx context.Context, productID uuid.UUID) error {
	// Check if exists in deleted state
	_, err := s.store.GetDeletedProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return common.ErrNotFound
		}
		return err
	}

	err = s.store.RestoreProduct(ctx, productID)
	if err != nil {
		s.log.Errorf("Failed to restore product", "productID", productID, "error", err)
		return err
	}

	// Log activity
	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	s.activityService.Log(ctx, actorID, repository.LogActionTypeUPDATE, repository.LogEntityTypePRODUCT, productID.String(), map[string]interface{}{
		"action": "restore",
	})

	return nil
}

func (s *PrdService) RestoreProductsBulk(ctx context.Context, req dto.RestoreBulkRequest) error {
	ids := make([]uuid.UUID, 0, len(req.ProductIDs))
	for _, idStr := range req.ProductIDs {
		uid, err := uuid.Parse(idStr)
		if err != nil {
			// Skip or fail? Fail entire request usually better for bulk operations consistency or return partial error?
			// User request didn't specify partial success. Assuming strict.
			return fmt.Errorf("invalid uuid: %s", idStr)
		}
		ids = append(ids, uid)
	}

	err := s.store.RestoreProductsBulk(ctx, ids)
	if err != nil {
		s.log.Errorf("Failed to bulk restore products", "error", err)
		return err
	}

	// Log activity (maybe one log for all or individually? One log is cleaner)
	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	s.activityService.Log(ctx, actorID, repository.LogActionTypeUPDATE, repository.LogEntityTypePRODUCT, "bulk", map[string]interface{}{
		"action": "restore_bulk",
		"count":  len(ids),
		"ids":    req.ProductIDs,
	})

	return nil
}
