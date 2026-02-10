package products

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"bytes"
	"errors"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IPrdHandler interface {
	CreateProductHandler(ctx *fiber.Ctx) error
	UploadProductImageHandler(ctx *fiber.Ctx) error
	ListProductsHandler(ctx *fiber.Ctx) error
	GetProductHandler(ctx *fiber.Ctx) error
	UpdateProductHandler(ctx *fiber.Ctx) error
	DeleteProductHandler(ctx *fiber.Ctx) error
	CreateProductOptionHandler(ctx *fiber.Ctx) error
	UploadProductOptionImageHandler(ctx *fiber.Ctx) error
	UpdateProductOptionHandler(ctx *fiber.Ctx) error
	DeleteProductOptionHandler(ctx *fiber.Ctx) error

	// Deleted Products
	ListDeletedProductsHandler(ctx *fiber.Ctx) error
	GetDeletedProductHandler(ctx *fiber.Ctx) error
	RestoreProductHandler(ctx *fiber.Ctx) error
	RestoreProductsBulkHandler(ctx *fiber.Ctx) error
	GetStockHistoryHandler(ctx *fiber.Ctx) error
}

func NewPrdHandler(prdService IPrdService, log logger.ILogger, validate validator.Validator) IPrdHandler {
	return &PrdHandler{
		prdService: prdService,
		log:        log,
		validate:   validate,
	}
}

type PrdHandler struct {
	prdService IPrdService
	log        logger.ILogger
	validate   validator.Validator
}

// DeleteProductOptionHandler
// @Summary Delete a product option
// @Description Delete a product option by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param option_id path string true "Option ID"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductOptionResponse} "Product option deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product or option not found"
// @Failure 500 {object} common.ErrorResponse "Failed to delete product option"
// @x-roles ["admin", "manager"]
// @Router /products/{product_id}/options/{option_id} [delete]

func (h *PrdHandler) DeleteProductOptionHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionIDStr := ctx.Params("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		h.log.Warn("Invalid option ID format", "error", err, "option_id", optionIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	err = h.prdService.DeleteProductOption(ctx.Context(), productID, optionID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product or option not found"})
		}
		h.log.Error("Failed to delete product option", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to delete product option"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product option deleted successfully",
	})
}

// UpdateProductOptionHandler
// @Summary Update a product option
// @Description Update a product option by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param option_id path string true "Option ID"
// @Param body body dto.UpdateProductOptionRequest true "Product option update request"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductOptionResponse} "Product option updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product or option not found"
// @Failure 500 {object} common.ErrorResponse "Failed to update product option"
// @x-roles ["admin", "manager"]
// @Router /products/{product_id}/options/{option_id} [patch]
func (h *PrdHandler) UpdateProductOptionHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionIDStr := ctx.Params("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		h.log.Warn("Invalid option ID format", "error", err, "option_id", optionIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	var req dto.UpdateProductOptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse product option update request body", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Product option update request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	optionResponse, err := h.prdService.UpdateProductOption(ctx.Context(), productID, optionID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product or option not found"})
		}
		h.log.Error("Failed to update product option", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update product option"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product option updated successfully",
		Data:    optionResponse,
	})
}

// UploadProductOptionImageHandler
// @Summary Upload product option image
// @Description Upload image for a specific product option
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param product_id path string true "Product ID"
// @Param option_id  path string true "Option ID"
// @Param image formData file true "Product image"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductOptionResponse} "Product option image uploaded successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid ID format or missing file"
// @Failure 404 {object} common.ErrorResponse "Product or option not found"
// @Failure 500 {object} common.ErrorResponse "Failed to upload product image"
// @x-roles ["admin", "manager"]
// @Router /products/{product_id}/options/{option_id}/image [post]
func (h *PrdHandler) UploadProductOptionImageHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionIDStr := ctx.Params("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		h.log.Warn("Invalid option ID format", "error", err, "option_id", optionIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		h.log.Warn("Image file is missing in form", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Image file is required in 'image' field"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Error("Failed to open uploaded file", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to process file"})
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			h.log.Error("Failed to close file", "error", err)
		}
	}(file)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		h.log.Error("Failed to read file into buffer", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to read file"})
	}

	optionResponse, err := h.prdService.UploadProductOptionImage(ctx.Context(), productID, optionID, buf.Bytes())
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product or option not found"})
		}
		h.log.Error("Failed to upload product option image", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to upload image"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product option image uploaded successfully",
		Data:    optionResponse,
	})
}

// CreateProductOptionHandler
// @Summary Create a product option
// @Description Create a product option for a parent product
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param body body dto.CreateProductOptionRequestStandalone true "Product option create request"
// @Success 201 {object} common.SuccessResponse{data=dto.ProductOptionResponse} "Product option created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Parent product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to create product option"
// @x-roles ["admin", "manager"]
// @Router /products/{product_id}/options [post]
func (h *PrdHandler) CreateProductOptionHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req dto.CreateProductOptionRequestStandalone
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse product option request body", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Product option request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	optionResponse, err := h.prdService.CreateProductOption(ctx.Context(), productID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Parent product not found"})
		}
		h.log.Error("Failed to create product option", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create product option"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Product option created successfully",
		Data:    optionResponse,
	})
}

// DeleteProductHandler
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} common.SuccessResponse "Product deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to delete product"
// @x-roles ["admin"]
// @Router /products/{id} [delete]
func (h *PrdHandler) DeleteProductHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	err = h.prdService.DeleteProduct(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to delete product", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to delete product"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// UpdateProductHandler
// @Summary Update a product
// @Description Update a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body dto.UpdateProductRequest true "Product update request"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductResponse} "Product updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to update product"
// @x-roles ["admin", "manager"]
// @Router /products/{id} [patch]
func (h *PrdHandler) UpdateProductHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req dto.UpdateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse request body for update", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Update request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	productResponse, err := h.prdService.UpdateProduct(ctx.Context(), productID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		if errors.Is(err, common.ErrCategoryNotFound) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Category not found"})
		}
		h.log.Error("Failed to update product", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update product"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product updated successfully",
		Data:    productResponse,
	})
}

// GetProductHandler
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductResponse} "Product retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve product"
// @x-roles ["admin", "manager", "cashier"]
// @Router /products/{id} [get]
func (h *PrdHandler) GetProductHandler(ctx *fiber.Ctx) error {
	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	productResponse, err := h.prdService.GetProductByID(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to get product by ID", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve product"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product retrieved successfully",
		Data:    productResponse,
	})
}

// ListProductsHandler
// @Summary List products
// @Description List products based on query parameters
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit the number of products returned"
// @Param search query string false "Search products by name"
// @Param category_id query int false "Search products by category ID"
// @Success 200 {object} common.SuccessResponse{data=dto.ListProductsResponse} "Products retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve products"
// @Router /products [get]
func (h *PrdHandler) ListProductsHandler(ctx *fiber.Ctx) error {
	var req dto.ListProductsRequest
	if err := ctx.QueryParser(&req); err != nil {
		h.log.Warn("Cannot parse query parameters", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
		})
	}

	h.log.Info("List products request received", "params", req)

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("List products request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	h.log.Info("List products request validated", "params", req)

	listResponse, err := h.prdService.ListProducts(ctx.Context(), req)
	if err != nil {
		h.log.Error("Failed to list products", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to retrieve products",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Products retrieved successfully",
		Data:    listResponse,
	})
}

// CreateProductHandler
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param body body dto.CreateProductRequest true "Product create request"
// @Success 201 {object} common.SuccessResponse{data=dto.ProductResponse} "Product created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request body"
// @Failure 409 {object} common.ErrorResponse "Product with same name already exists"
// @Failure 500 {object} common.ErrorResponse "Failed to create product"
// @x-roles ["admin", "manager"]
// @Router /products [post]
func (h *PrdHandler) CreateProductHandler(ctx *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse request body", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	productResponse, err := h.prdService.CreateProduct(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to create product",
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Product created successfully",
		Data:    productResponse,
	})
}

// UploadProductImageHandler
// @Summary Upload an image for a product
// @Description Upload an image for a product by ID
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param image formData file true "Image file"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductResponse} "Product image uploaded successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format or image file is missing"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to upload image"
// @x-roles ["admin", "manager"]
// @Router /products/{id}/image [post]
func (h *PrdHandler) UploadProductImageHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		h.log.Warn("Image file is missing or invalid in form", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Image file is required in 'image' field"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Error("Failed to open uploaded file", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to process file"})
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			h.log.Error("Failed to close file", "error", err)
		}
	}(file)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		h.log.Error("Failed to read file into buffer", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to read file"})
	}
	fileBytes := buf.Bytes()

	productResponse, err := h.prdService.UploadProductImage(ctx.Context(), productID, fileBytes)
	if err != nil {

		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to upload image",
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product image uploaded successfully",
		Data:    productResponse,
	})
}

// ListDeletedProductsHandler
// @Summary List deleted products
// @Description List deleted products with pagination and filtering
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit the number of products returned"
// @Param search query string false "Search products by name"
// @Param category_id query int false "Search products by category ID"
// @Success 200 {object} common.SuccessResponse{data=dto.ListProductsResponse} "Deleted products retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve deleted products"
// @x-roles ["admin"]
// @Router /products/trash [get]
func (h *PrdHandler) ListDeletedProductsHandler(ctx *fiber.Ctx) error {
	var req dto.ListProductsRequest
	if err := ctx.QueryParser(&req); err != nil {
		h.log.Warn("Cannot parse query parameters", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
		})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("List deleted products request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	listResponse, err := h.prdService.ListDeletedProducts(ctx.Context(), req)
	if err != nil {
		h.log.Error("Failed to list deleted products", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to retrieve deleted products",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Deleted products retrieved successfully",
		Data:    listResponse,
	})
}

// GetDeletedProductHandler
// @Summary Get a deleted product
// @Description Get a deleted product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} common.SuccessResponse{data=dto.ProductResponse} "Deleted product retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve deleted product"
// @x-roles ["admin"]
// @Router /products/trash/{id} [get]
func (h *PrdHandler) GetDeletedProductHandler(ctx *fiber.Ctx) error {
	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	productResponse, err := h.prdService.GetDeletedProduct(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to get deleted product by ID", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve deleted product"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Deleted product retrieved successfully",
		Data:    productResponse,
	})
}

// RestoreProductHandler
// @Summary Restore a deleted product
// @Description Restore a deleted product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} common.SuccessResponse "Product restored successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to restore product"
// @x-roles ["admin"]
// @Router /products/trash/{id}/restore [post]
func (h *PrdHandler) RestoreProductHandler(ctx *fiber.Ctx) error {
	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	err = h.prdService.RestoreProduct(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to restore product", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to restore product"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product restored successfully",
	})
}

// RestoreProductsBulkHandler
// @Summary Bulk restore deleted products
// @Description Restore multiple deleted products by IDs
// @Tags Products
// @Accept json
// @Produce json
// @Param body body dto.RestoreBulkRequest true "Bulk restore request"
// @Success 200 {object} common.SuccessResponse "Products restored successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request body"
// @Failure 500 {object} common.ErrorResponse "Failed to restore products"
// @x-roles ["admin"]
// @Router /products/trash/restore-bulk [post]
func (h *PrdHandler) RestoreProductsBulkHandler(ctx *fiber.Ctx) error {
	var req dto.RestoreBulkRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse request body", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Bulk restore request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	err := h.prdService.RestoreProductsBulk(ctx.Context(), req)
	if err != nil {
		h.log.Error("Failed to bulk restore products", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to restore products",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Products restored successfully",
	})
}

// GetStockHistoryHandler
// @Summary Get stock history for a product
// @Description Get stock history for a product by ID with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} common.SuccessResponse{data=dto.PagedStockHistoryResponse} "Stock history retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid product ID or query parameters"
// @Failure 404 {object} common.ErrorResponse "Product not found"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve stock history"
// @x-roles ["admin", "manager", "cashier"]
// @Router /products/{id}/stock-history [get]
func (h *PrdHandler) GetStockHistoryHandler(ctx *fiber.Ctx) error {
	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req dto.ListStockHistoryRequest
	if err := ctx.QueryParser(&req); err != nil {
		h.log.Warn("Cannot parse query parameters", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
		})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Stock history request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	historyResponse, err := h.prdService.GetStockHistory(ctx.Context(), productID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to get stock history", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve stock history"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Stock history retrieved successfully",
		Data:    historyResponse,
	})
}
