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
}

func NewPrdHandler(prdService IPrdService, log *logger.Logger, validate validator.Validator) IPrdHandler {
	return &PrdHandler{
		prdService: prdService,
		log:        log,
		validate:   validate,
	}
}

type PrdHandler struct {
	prdService IPrdService
	log        *logger.Logger
	validate   validator.Validator
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
