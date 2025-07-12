package products

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/pagination"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"sync"
)

type IPrdService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	UploadProductImage(ctx context.Context, productID uuid.UUID, data []byte) (*ProductResponse, error)
	ListProducts(ctx context.Context, req ListProductsRequest) (*ListProductsResponse, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*ProductResponse, error)
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

func (s *PrdService) GetProductByID(ctx context.Context, productID uuid.UUID) (*ProductResponse, error) {
	fullProduct, err := s.store.GetProductWithOptions(ctx, productID)
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

func (s *PrdService) buildProductResponse(ctx context.Context, fullProduct repository.GetProductWithOptionsRow) (*ProductResponse, error) {
	var optionsResponse []ProductOptionResponse
	if fullProduct.Options != nil {
		if optionsJSON, ok := fullProduct.Options.([]byte); ok {
			var options []repository.ProductOption
			if err := json.Unmarshal(optionsJSON, &options); err != nil {
				s.log.Error("Failed to unmarshal product options JSON", "error", err)
				return nil, fmt.Errorf("could not parse product options")
			}
			for _, opt := range options {
				var additionalPrice float64
				_ = opt.AdditionalPrice.Scan(&additionalPrice)
				optionsResponse = append(optionsResponse, ProductOptionResponse{
					ID:              opt.ID,
					Name:            opt.Name,
					AdditionalPrice: additionalPrice,
					ImageURL:        opt.ImageUrl,
				})
			}
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
				url, err := s.prdRepo.PrdImageLink(ctx, opt.ID.String(), *opt.ImageURL)
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

	var productPrice float64
	_ = fullProduct.Price.Scan(&productPrice)

	return &ProductResponse{
		ID:         fullProduct.ID,
		Name:       fullProduct.Name,
		CategoryID: fullProduct.CategoryID,
		ImageURL:   fullProduct.ImageUrl,
		Price:      productPrice,
		Stock:      fullProduct.Stock,
		CreatedAt:  fullProduct.CreatedAt.Time,
		UpdatedAt:  fullProduct.UpdatedAt.Time,
		Options:    optionsResponse,
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
		var price float64
		_ = p.Price.Scan(&price)
		productsResponse = append(productsResponse, ProductListResponse{
			ID:           p.ID,
			Name:         p.Name,
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
		var priceNumeric pgtype.Numeric
		if err = priceNumeric.Scan(req.Price); err != nil {
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
			var additionalPriceNumeric pgtype.Numeric
			if err = additionalPriceNumeric.Scan(opt.AdditionalPrice); err != nil {
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

	const maxFileSize = 2 * 1024 * 1024 // 2MB
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
		ImageUrl: &imageUrl,
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
		var additionalPrice float64
		_ = opt.AdditionalPrice.Scan(&additionalPrice)
		optionsResponse = append(optionsResponse, ProductOptionResponse{
			ID:              opt.ID,
			Name:            opt.Name,
			AdditionalPrice: additionalPrice,
			ImageURL:        opt.ImageUrl,
		})
	}

	var productPrice float64
	_ = product.Price.Scan(&productPrice)

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
