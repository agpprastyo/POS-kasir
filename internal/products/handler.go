package products

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
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

// DeleteProductOptionHandler menangani penghapusan (soft delete) varian produk.
func (h *PrdHandler) DeleteProductOptionHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warnf("invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionIDStr := ctx.Params("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		h.log.Warnf("Invalid option ID format", "error", err, "option_id", optionIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	err = h.prdService.DeleteProductOption(ctx.Context(), productID, optionID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product or option not found"})
		}
		h.log.Errorf("Failed to delete product option", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to delete product option"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product option deleted successfully",
	})
}

func (h *PrdHandler) UpdateProductOptionHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warnf("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionIDStr := ctx.Params("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		h.log.Warnf("Invalid option ID format", "error", err, "option_id", optionIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	// 2. Parse request body ke DTO
	var req UpdateProductOptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warnf("Cannot parse product option update request body", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Product option update request validation failed", "error", err)
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
		h.log.Errorf("Failed to update product option", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update product option"})
	}

	// 5. Kirim respons sukses
	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product option updated successfully",
		Data:    optionResponse,
	})
}

func (h *PrdHandler) UploadProductOptionImageHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warnf("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionIDStr := ctx.Params("option_id")
	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		h.log.Warnf("Invalid option ID format", "error", err, "option_id", optionIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		h.log.Warnf("Image file is missing in form", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Image file is required in 'image' field"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Errorf("Failed to open uploaded file", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to process file"})
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			h.log.Errorf("Failed to close file", "error", err)
		}
	}(file)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		h.log.Errorf("Failed to read file into buffer", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to read file"})
	}

	optionResponse, err := h.prdService.UploadProductOptionImage(ctx.Context(), productID, optionID, buf.Bytes())
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product or option not found"})
		}
		h.log.Errorf("Failed to upload product option image", "error", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to upload image"})
	}

	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product option image uploaded successfully",
		Data:    optionResponse,
	})
}

// CreateProductOptionHandler menangani pembuatan varian baru untuk produk yang sudah ada.
func (h *PrdHandler) CreateProductOptionHandler(ctx *fiber.Ctx) error {

	productIDStr := ctx.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warnf("Invalid product ID format", "error", err, "product_id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req CreateProductOptionRequestStandalone
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warnf("Cannot parse product option request body", "error", err)
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
		h.log.Errorf("Failed to create product option", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create product option"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Product option created successfully",
		Data:    optionResponse,
	})
}

// DeleteProductHandler menangani request untuk menghapus (soft delete) produk.
func (h *PrdHandler) DeleteProductHandler(ctx *fiber.Ctx) error {
	// 1. Ambil ID produk dari parameter URL
	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	// 2. Panggil service untuk melakukan soft delete
	err = h.prdService.DeleteProduct(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to delete product", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to delete product"})
	}

	// 3. Kirim respons sukses
	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// UpdateProductHandler menangani request untuk memperbarui data produk.
func (h *PrdHandler) UpdateProductHandler(ctx *fiber.Ctx) error {
	// 1. Ambil ID produk dari parameter URL
	productIDStr := ctx.Params("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", productIDStr)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	// 2. Parse request body ke DTO
	var req UpdateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse request body for update", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	// 3. Validasi DTO
	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Update request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// 4. Panggil service untuk melakukan update
	productResponse, err := h.prdService.UpdateProduct(ctx.Context(), productID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Product not found"})
		}
		h.log.Error("Failed to update product", "error", err, "productID", productID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update product"})
	}

	// 5. Kirim respons sukses
	return ctx.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Product updated successfully",
		Data:    productResponse,
	})
}

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

func (h *PrdHandler) ListProductsHandler(ctx *fiber.Ctx) error {
	var req ListProductsRequest
	if err := ctx.QueryParser(&req); err != nil {
		h.log.Warn("Cannot parse query parameters", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
		})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("List products request validation failed", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// parse manual

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
func (h *PrdHandler) CreateProductHandler(ctx *fiber.Ctx) error {
	var req CreateProductRequest
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
