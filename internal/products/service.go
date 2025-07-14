package products

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/pagination"
	"POS-kasir/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sync"
)

type IPrdService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	UploadProductImage(ctx context.Context, productID uuid.UUID, data []byte) (*ProductResponse, error)
	ListProducts(ctx context.Context, req ListProductsRequest) (*ListProductsResponse, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	CreateProductOption(ctx context.Context, productID uuid.UUID, req CreateProductOptionRequestStandalone) (*ProductOptionResponse, error)
	UploadProductOptionImage(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, data []byte) (*ProductOptionResponse, error)
	UpdateProductOption(ctx context.Context, productID, optionID uuid.UUID, req UpdateProductOptionRequest) (*ProductOptionResponse, error)
	DeleteProductOption(ctx context.Context, productID, optionID uuid.UUID) error
}

type PrdService struct {
	log             *logger.Logger
	store           repository.Store
	prdRepo         IPrdRepo
	activityService activitylog.Service
}

func NewPrdService(store repository.Store, log *logger.Logger, prdRepo IPrdRepo, activityService activitylog.Service) IPrdService {
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
			s.log.Warn("Product option not found or does not belong to the product", "optionID", optionID, "productID", productID)
			return common.ErrNotFound
		}
		s.log.Error("Failed to get product option before deletion", "error", err)
		return err
	}

	err = s.store.SoftDeleteProductOption(ctx, optionID)
	if err != nil {
		s.log.Error("Failed to soft delete product option in repository", "error", err, "optionID", optionID)
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

	s.log.Info("Product option soft deleted successfully", "optionID", optionID)
	return nil
}

func (s *PrdService) UpdateProductOption(ctx context.Context, productID, optionID uuid.UUID, req UpdateProductOptionRequest) (*ProductOptionResponse, error) {
	_, err := s.store.GetProductOption(ctx, repository.GetProductOptionParams{ID: optionID, ProductID: productID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Product option not found or does not belong to the product", "optionID", optionID, "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get product option before update", "error", err)
		return nil, err
	}

	updateParams := repository.UpdateProductOptionParams{
		ID:   optionID,
		Name: req.Name,
	}

	if req.AdditionalPrice != nil {
		priceNumeric, err := utils.Float64ToNumeric(*req.AdditionalPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to convert additional price for update: %w", err)
		}
		updateParams.AdditionalPrice = priceNumeric
	}

	updatedOption, err := s.store.UpdateProductOption(ctx, updateParams)
	if err != nil {
		s.log.Error("Failed to update product option in repository", "error", err, "optionID", optionID)
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

	additionalPrice := utils.NumericToFloat64(updatedOption.AdditionalPrice)

	if updatedOption.ImageUrl != nil {
		publicUrl, err := s.prdRepo.PrdImageLink(ctx, updatedOption.ID.String(), *updatedOption.ImageUrl)
		if err != nil {
			s.log.Warn("Failed to get public URL for updated option image", "error", err)
			publicUrl = *updatedOption.ImageUrl
		}
		updatedOption.ImageUrl = &publicUrl
	} else {
		updatedOption.ImageUrl = nil // Set to nil if no image URL is present
	}

	return &ProductOptionResponse{
		ID:              updatedOption.ID,
		Name:            updatedOption.Name,
		AdditionalPrice: additionalPrice,
		ImageURL:        updatedOption.ImageUrl,
	}, nil
}

func (s *PrdService) UploadProductOptionImage(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, data []byte) (*ProductOptionResponse, error) {
	_, err := s.store.GetProductOption(ctx, repository.GetProductOptionParams{
		ID:        optionID,
		ProductID: productID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Product option not found or does not belong to the product", "optionID", optionID, "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get product option before image upload", "error", err)
		return nil, err
	}

	const maxFileSize = 2 * 1024 * 1024
	if len(data) > maxFileSize {
		return nil, fmt.Errorf("file size exceeds the limit of 2MB")
	}

	filename := fmt.Sprintf("product_options/%s.jpg", optionID.String())

	imageUrl, err := s.prdRepo.UploadImageToMinio(ctx, filename, data)
	if err != nil {
		s.log.Error("Failed to upload option image to Minio", "error", err)
		return nil, fmt.Errorf("could not upload image to storage")
	}

	updateParams := repository.UpdateProductOptionParams{
		ID:       optionID,
		ImageUrl: &filename,
	}
	updatedOption, err := s.store.UpdateProductOption(ctx, updateParams)
	if err != nil {
		s.log.Error("Failed to update product option with image URL", "error", err)
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
	_ = updatedOption.AdditionalPrice.Scan(&additionalPrice)

	publicUrl, err := s.prdRepo.PrdImageLink(ctx, updatedOption.ID.String(), *updatedOption.ImageUrl)
	if err != nil {
		s.log.Warn("Failed to get public URL for newly uploaded option image", "error", err)
		publicUrl = *updatedOption.ImageUrl
	}

	return &ProductOptionResponse{
		ID:              updatedOption.ID,
		Name:            updatedOption.Name,
		AdditionalPrice: additionalPrice,
		ImageURL:        &publicUrl,
	}, nil
}
func (s *PrdService) CreateProductOption(ctx context.Context, productID uuid.UUID, req CreateProductOptionRequestStandalone) (*ProductOptionResponse, error) {
	_, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Parent product not found for new option", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get parent product for new option", "error", err)
		return nil, err
	}

	additionalPriceNumeric, err := utils.Float64ToNumeric(req.AdditionalPrice)
	if err != nil {
		s.log.Error("Failed to convert additional price to numeric", "error", err)
		return nil, fmt.Errorf("failed to convert additional price: %w", err)
	}
	params := repository.CreateProductOptionParams{
		ProductID:       productID,
		Name:            req.Name,
		AdditionalPrice: additionalPriceNumeric,
	}
	newOption, err := s.store.CreateProductOption(ctx, params)
	if err != nil {
		s.log.Error("Failed to create product option in repository", "error", err)
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

	var additionalPrice = utils.NumericToFloat64(newOption.AdditionalPrice)

	return &ProductOptionResponse{
		ID:              newOption.ID,
		Name:            newOption.Name,
		AdditionalPrice: additionalPrice,
		ImageURL:        newOption.ImageUrl,
	}, nil
}
func (s *PrdService) DeleteProduct(ctx context.Context, productID uuid.UUID) error {

	product, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Product not found for deletion", "productID", productID)
			return common.ErrNotFound
		}
		s.log.Error("Failed to get product before deletion", "error", err)
		return err
	}

	err = s.store.SoftDeleteProduct(ctx, productID)
	if err != nil {
		s.log.Error("Failed to soft delete product in repository", "error", err, "productID", productID)
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

	s.log.Info("Product soft deleted successfully", "productID", productID)
	return nil
}

func (s *PrdService) GetProductByID(ctx context.Context, productID uuid.UUID) (*ProductResponse, error) {
	fullProduct, err := s.store.GetProductWithOptions(ctx, productID)
	s.log.Infof("Full product data: %+v", fullProduct)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Product not found by ID", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get product from repository", "error", err, "productID", productID)
		return nil, err
	}

	return s.buildProductResponse(ctx, fullProduct)
}

func (s *PrdService) UpdateProduct(ctx context.Context, productID uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	// 1. Pastikan produk ada sebelum update
	_, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Product not found for update", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get product before update", "error", err)
		return nil, err
	}

	updateParams := repository.UpdateProductParams{
		ID:         productID,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Stock:      req.Stock,
	}

	if req.Price != nil {
		priceNumeric, err := utils.Float64ToNumeric(*req.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price for update: %w", err)
		}
		updateParams.Price = priceNumeric
	}

	_, err = s.store.UpdateProduct(ctx, updateParams)
	if err != nil {
		s.log.Error("Failed to update product in repository", "error", err, "productID", productID)
		return nil, err
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
func (s *PrdService) buildProductResponse(ctx context.Context, fullProduct repository.GetProductWithOptionsRow) (*ProductResponse, error) {
	var optionsResponse []ProductOptionResponse

	// --- PERBAIKAN LOGIKA UNMARSHAL ---
	if fullProduct.Options != nil {
		// 1. Marshal kembali interface{} menjadi []byte JSON.
		// Ini adalah cara andal untuk menangani hasil dari json_agg.
		optionsJSON, err := json.Marshal(fullProduct.Options)
		if err != nil {
			s.log.Error("Failed to re-marshal product options interface", "error", err)
			return nil, fmt.Errorf("could not process product options")
		}

		// 2. Unmarshal []byte JSON tersebut ke dalam slice struct yang benar.
		var options []repository.ProductOption
		if err := json.Unmarshal(optionsJSON, &options); err != nil {
			// Jika unmarshal gagal, mungkin karena datanya kosong ('[]')
			if string(optionsJSON) != "[]" {
				s.log.Error("Failed to unmarshal product options JSON", "error", err)
				return nil, fmt.Errorf("could not parse product options")
			}
		}

		// 3. Lakukan loop pada slice struct yang sudah benar.
		for _, opt := range options {
			additionalPrice := utils.NumericToFloat64(opt.AdditionalPrice)
			optionsResponse = append(optionsResponse, ProductOptionResponse{
				ID:              opt.ID,
				Name:            opt.Name,
				AdditionalPrice: additionalPrice,
				ImageURL:        opt.ImageUrl,
			})
		}
	}
	// --- AKHIR PERBAIKAN ---

	var wg sync.WaitGroup
	errChan := make(chan error, len(optionsResponse)+1)

	// Ambil URL untuk gambar produk utama
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

	// Ambil URL untuk setiap gambar varian
	for i := range optionsResponse {
		// Gunakan variabel lokal di dalam loop untuk goroutine
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

	// Cek apakah ada error dari goroutine
	for err := range errChan {
		if err != nil {
			return nil, err // Kembalikan error pertama yang ditemui
		}
	}

	productPrice := utils.NumericToFloat64(fullProduct.Price)

	return &ProductResponse{
		ID:         fullProduct.ID,
		Name:       fullProduct.Name,
		CategoryID: fullProduct.CategoryID,

		ImageURL:  fullProduct.ImageUrl,
		Price:     productPrice,
		Stock:     fullProduct.Stock,
		CreatedAt: fullProduct.CreatedAt.Time,
		UpdatedAt: fullProduct.UpdatedAt.Time,
		Options:   optionsResponse,
	}, nil
}

func (s *PrdService) ListProducts(ctx context.Context, req ListProductsRequest) (*ListProductsResponse, error) {

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
		s.log.Error("Failed to list products from repository", "error", listErr)
		return nil, listErr
	}
	if countErr != nil {
		s.log.Error("Failed to count products from repository", "error", countErr)
		return nil, countErr
	}

	var productsResponse []ProductListResponse
	for _, p := range products {
		price := utils.NumericToFloat64(p.Price)
		productsResponse = append(productsResponse, ProductListResponse{
			ID:           p.ID,
			Name:         p.Name,
			CategoryID:   p.CategoryID,
			CategoryName: p.CategoryName,
			ImageURL:     p.ImageUrl,
			Price:        price,
			Stock:        p.Stock,
		})
	}

	response := &ListProductsResponse{
		Products: productsResponse,
		Pagination: pagination.BuildPagination(
			page,
			int(totalData),
			limit,
		),
	}

	return response, nil
}

func (s *PrdService) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	var newProduct repository.Product
	var createdOptions []repository.ProductOption

	txFunc := func(qtx *repository.Queries) error {
		var err error

		priceNumeric, err := utils.Float64ToNumeric(req.Price)
		if err != nil {
			return fmt.Errorf("failed to scan price: %w", err)
		}

		productParams := repository.CreateProductParams{
			Name:       req.Name,
			CategoryID: &req.CategoryID,
			Price:      priceNumeric,
			Stock:      req.Stock,
		}
		newProduct, err = qtx.CreateProduct(ctx, productParams)
		if err != nil {
			s.log.Error("Failed to create product in transaction", "error", err)
			return err
		}

		for _, opt := range req.Options {
			additionalPriceNumeric, err := utils.Float64ToNumeric(opt.AdditionalPrice)
			if err != nil {
				return fmt.Errorf("failed to scan additional price for option %s: %w", opt.Name, err)
			}

			optionParams := repository.CreateProductOptionParams{
				ProductID:       newProduct.ID,
				Name:            opt.Name,
				AdditionalPrice: additionalPriceNumeric,
			}
			createdOpt, err := qtx.CreateProductOption(ctx, optionParams)
			if err != nil {
				s.log.Error("Failed to create product option in transaction", "error", err)
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

func (s *PrdService) UploadProductImage(ctx context.Context, productID uuid.UUID, data []byte) (*ProductResponse, error) {
	_, err := s.store.GetProductWithOptions(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Product not found for image upload", "productID", productID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get product for image upload", "error", err)
		return nil, err
	}

	const maxFileSize = 5 * 1024 * 1024 // 2MB
	if len(data) > maxFileSize {
		return nil, fmt.Errorf("file size exceeds the limit of 2MB")
	}

	filename := fmt.Sprintf("products/%s.jpg", productID.String())

	imageUrl, err := s.prdRepo.UploadImageToMinio(ctx, filename, data)
	if err != nil {
		s.log.Error("Failed to upload image to Minio", "error", err)
		return nil, fmt.Errorf("could not upload image to storage")
	}

	updateParams := repository.UpdateProductParams{
		ID:       productID,
		ImageUrl: &filename,
	}
	_, err = s.store.UpdateProduct(ctx, updateParams)
	if err != nil {
		s.log.Error("Failed to update product with image URL", "error", err)
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
		s.log.Error("Failed to fetch full product after image upload", "error", err)
		return nil, err
	}

	return s.buildProductResponse(ctx, fullProduct)
}

func (s *PrdService) buildProductResponseFromData(product repository.Product, options []repository.ProductOption) (*ProductResponse, error) {
	var optionsResponse []ProductOptionResponse
	for _, opt := range options {
		var additionalPrice = utils.NumericToFloat64(opt.AdditionalPrice)
		optionsResponse = append(optionsResponse, ProductOptionResponse{
			ID:              opt.ID,
			Name:            opt.Name,
			AdditionalPrice: additionalPrice,
			ImageURL:        opt.ImageUrl,
		})
	}

	productPrice := utils.NumericToFloat64(product.Price)

	s.log.Infof("product price before assign: %+v", product.Price)
	s.log.Infof("product price after assign: %+v", productPrice)

	return &ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryID,
		ImageURL:   product.ImageUrl,
		Price:      productPrice,
		Stock:      product.Stock,
		CreatedAt:  product.CreatedAt.Time,
		UpdatedAt:  product.UpdatedAt.Time,
		Options:    optionsResponse,
	}, nil
}
