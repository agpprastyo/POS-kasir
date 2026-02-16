package products

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"bytes"
	"errors"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type IPrdHandler interface {
	CreateProductHandler(ctx fiber.Ctx) error
	UploadProductImageHandler(ctx fiber.Ctx) error
	ListProductsHandler(ctx fiber.Ctx) error
	GetProductHandler(ctx fiber.Ctx) error
	UpdateProductHandler(ctx fiber.Ctx) error
	DeleteProductHandler(ctx fiber.Ctx) error
	CreateProductOptionHandler(ctx fiber.Ctx) error
	UploadProductOptionImageHandler(ctx fiber.Ctx) error
	UpdateProductOptionHandler(ctx fiber.Ctx) error
	DeleteProductOptionHandler(ctx fiber.Ctx) error

	// Deleted Products
	ListDeletedProductsHandler(ctx fiber.Ctx) error
	GetDeletedProductHandler(ctx fiber.Ctx) error
	RestoreProductHandler(ctx fiber.Ctx) error
	RestoreProductsBulkHandler(ctx fiber.Ctx) error
	GetStockHistoryHandler(ctx fiber.Ctx) error
}

func NewPrdHandler(prdService IPrdService, log logger.ILogger) IPrdHandler {
	return &PrdHandler{
		prdService: prdService,
		log:        log,
	}
}

type PrdHandler struct {
	prdService IPrdService
	log        logger.ILogger
}

// DeleteProductOptionHandler deletes a product option
// @Summary      Delete a product option
// @Description  Delete a specific product option by its ID (Roles: admin, manager)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id path string true "Product ID" Format(uuid)
// @Param        option_id  path string true "Option ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Product option deleted successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format"
// @Failure      404 {object} common.ErrorResponse "Product or option not found"
// @Failure      500 {object} common.ErrorResponse "Failed to delete product option"
// @x-roles      ["admin", "manager"]
// @Router       /products/{product_id}/options/{option_id} [delete]

func (h *PrdHandler) DeleteProductOptionHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("product_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", ctx.Params("product_id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	optionID, err := fiber.Convert(ctx.Params("option_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid option ID format", "error", err, "option_id", ctx.Params("option_id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	err = h.prdService.DeleteProductOption(ctx.RequestCtx(), productID, optionID)
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

// UpdateProductOptionHandler updates a product option
// @Summary      Update a product option
// @Description  Update details of a specific product option by its ID (Roles: admin, manager)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id path string true "Product ID" Format(uuid)
// @Param        option_id  path string true "Option ID" Format(uuid)
// @Param        body       body UpdateProductOptionRequest true "Product option update request"
// @Success      200 {object} common.SuccessResponse{data=ProductOptionResponse} "Product option updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Product or option not found"
// @Failure      500 {object} common.ErrorResponse "Failed to update product option"
// @x-roles      ["admin", "manager"]
// @Router       /products/{product_id}/options/{option_id} [patch]
func (h *PrdHandler) UpdateProductOptionHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("product_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", ctx.Params("product_id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionID, err := fiber.Convert(ctx.Params("option_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid option ID format", "error", err, "option_id", ctx.Params("option_id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid option ID format"})
	}

	var req UpdateProductOptionRequest
	if err := ctx.Bind().Body(&req); err != nil {
		h.log.Warn("Product option update request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body", Error: err.Error()})
	}

	optionResponse, err := h.prdService.UpdateProductOption(ctx.RequestCtx(), productID, optionID, req)
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

// UploadProductOptionImageHandler uploads an image for a product option
// @Summary      Upload product option image
// @Description  Upload image for a specific product option (Roles: admin, manager)
// @Tags         Products
// @Accept       multipart/form-data
// @Produce      json
// @Param        product_id path     string true "Product ID" Format(uuid)
// @Param        option_id  path     string true "Option ID" Format(uuid)
// @Param        image      formData file   true "Product option image"
// @Success      200 {object} common.SuccessResponse{data=ProductOptionResponse} "Product option image uploaded successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format or missing file"
// @Failure      404 {object} common.ErrorResponse "Product or option not found"
// @Failure      500 {object} common.ErrorResponse "Failed to upload product option image"
// @x-roles      ["admin", "manager"]
// @Router       /products/{product_id}/options/{option_id}/image [post]
func (h *PrdHandler) UploadProductOptionImageHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("product_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", ctx.Params("product_id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}
	optionID, err := fiber.Convert(ctx.Params("option_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid option ID format", "error", err, "option_id", ctx.Params("option_id"))
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

	optionResponse, err := h.prdService.UploadProductOptionImage(ctx.RequestCtx(), productID, optionID, buf.Bytes())
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

// CreateProductOptionHandler creates a new product option
// @Summary      Create a product option
// @Description  Create a new product option for a parent product (Roles: admin, manager)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id path string true "Product ID" Format(uuid)
// @Param        body       body CreateProductOptionRequestStandalone true "Product option create request"
// @Success      201 {object} common.SuccessResponse{data=ProductOptionResponse} "Product option created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid product ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Parent product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to create product option"
// @x-roles      ["admin", "manager"]
// @Router       /products/{product_id}/options [post]
func (h *PrdHandler) CreateProductOptionHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("product_id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "product_id", ctx.Params("product_id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req CreateProductOptionRequestStandalone
	if err := ctx.Bind().Body(&req); err != nil {
		h.log.Warn("Product option request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body", Error: err.Error()})
	}

	optionResponse, err := h.prdService.CreateProductOption(ctx.RequestCtx(), productID, req)
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

// DeleteProductHandler deletes a product
// @Summary      Delete a product
// @Description  Delete a product by its ID (Roles: admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Product deleted successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to delete product"
// @x-roles      ["admin"]
// @Router       /products/{id} [delete]
func (h *PrdHandler) DeleteProductHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	err = h.prdService.DeleteProduct(ctx.RequestCtx(), productID)
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

// UpdateProductHandler updates a product
// @Summary      Update a product
// @Description  Update details of a specific product by its ID (Roles: admin, manager)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path string true "Product ID" Format(uuid)
// @Param        body body UpdateProductRequest true "Product update request"
// @Success      200 {object} common.SuccessResponse{data=ProductResponse} "Product updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to update product"
// @x-roles      ["admin", "manager"]
// @Router       /products/{id} [patch]
func (h *PrdHandler) UpdateProductHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req UpdateProductRequest
	if err := ctx.Bind().Body(&req); err != nil {
		h.log.Warn("Update request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body", Error: err.Error()})
	}

	productResponse, err := h.prdService.UpdateProduct(ctx.RequestCtx(), productID, req)
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

// GetProductHandler gets a product by ID
// @Summary      Get a product by ID
// @Description  Retrieve detailed information of a specific product by its ID (Roles: admin, manager, cashier)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=ProductResponse} "Product retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve product"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /products/{id} [get]
func (h *PrdHandler) GetProductHandler(ctx fiber.Ctx) error {
	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	productResponse, err := h.prdService.GetProductByID(ctx.RequestCtx(), productID)
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

// ListProductsHandler lists all products with filtering and pagination
// @Summary      List products
// @Description  Get a list of products with filtering by category and search term (Roles: authenticated)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page        query int    false "Page number"
// @Param        limit       query int    false "Limit the number of products returned"
// @Param        search      query string false "Search products by name"
// @Param        category_id query int    false "Search products by category ID"
// @Success      200 {object} common.SuccessResponse{data=ListProductsResponse} "Products retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve products"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /products [get]
func (h *PrdHandler) ListProductsHandler(ctx fiber.Ctx) error {
	var req ListProductsRequest
	if err := ctx.Bind().Query(&req); err != nil {
		h.log.Warn("List products request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	h.log.Infof("List products request received: %+v", req)

	listResponse, err := h.prdService.ListProducts(ctx.RequestCtx(), req)
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

// CreateProductHandler creates a new product
// @Summary      Create a new product
// @Description  Create a new product with multiple options (Roles: admin, manager)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        body body CreateProductRequest true "Product create request"
// @Success      201 {object} common.SuccessResponse{data=ProductResponse} "Product created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body"
// @Failure      409 {object} common.ErrorResponse "Product with same name already exists"
// @Failure      500 {object} common.ErrorResponse "Failed to create product"
// @x-roles      ["admin", "manager"]
// @Router       /products [post]
func (h *PrdHandler) CreateProductHandler(ctx fiber.Ctx) error {
	var req CreateProductRequest
	if err := ctx.Bind().Body(&req); err != nil {
		h.log.Warn("Request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	productResponse, err := h.prdService.CreateProduct(ctx.RequestCtx(), req)
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

// UploadProductImageHandler uploads an image for a product
// @Summary      Upload an image for a product
// @Description  Upload an image for a product by ID (Roles: admin, manager)
// @Tags         Products
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path     string true "Product ID" Format(uuid)
// @Param        image formData file   true "Image file"
// @Success      200 {object} common.SuccessResponse{data=ProductResponse} "Product image uploaded successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format or image file is missing"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to upload image"
// @x-roles      ["admin", "manager"]
// @Router       /products/{id}/image [post]
func (h *PrdHandler) UploadProductImageHandler(ctx fiber.Ctx) error {

	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
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

	productResponse, err := h.prdService.UploadProductImage(ctx.RequestCtx(), productID, fileBytes)
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

// ListDeletedProductsHandler lists all deleted products
// @Summary      List deleted products
// @Description  Get a list of deleted products with pagination and filtering (Roles: admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        page        query int    false "Page number"
// @Param        limit       query int    false "Limit the number of products returned"
// @Param        search      query string false "Search products by name"
// @Param        category_id query int    false "Search products by category ID"
// @Success      200 {object} common.SuccessResponse{data=ListProductsResponse} "Deleted products retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve deleted products"
// @x-roles      ["admin"]
// @Router       /products/trash [get]
func (h *PrdHandler) ListDeletedProductsHandler(ctx fiber.Ctx) error {
	var req ListProductsRequest
	if err := ctx.Bind().Query(&req); err != nil {
		h.log.Warn("List deleted products request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	listResponse, err := h.prdService.ListDeletedProducts(ctx.RequestCtx(), req)
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

// GetDeletedProductHandler gets a deleted product by ID
// @Summary      Get a deleted product
// @Description  Retrieve detailed information of a specific deleted product by its ID (Roles: admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=ProductResponse} "Deleted product retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve deleted product"
// @x-roles      ["admin"]
// @Router       /products/trash/{id} [get]
func (h *PrdHandler) GetDeletedProductHandler(ctx fiber.Ctx) error {
	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	productResponse, err := h.prdService.GetDeletedProduct(ctx.RequestCtx(), productID)
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

// RestoreProductHandler restores a deleted product
// @Summary      Restore a deleted product
// @Description  Restore a specific deleted product by its ID (Roles: admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Product restored successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid product ID format"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to restore product"
// @x-roles      ["admin"]
// @Router       /products/trash/{id}/restore [post]
func (h *PrdHandler) RestoreProductHandler(ctx fiber.Ctx) error {
	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	err = h.prdService.RestoreProduct(ctx.RequestCtx(), productID)
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

// RestoreProductsBulkHandler restores multiple deleted products
// @Summary      Bulk restore deleted products
// @Description  Restore multiple deleted products by their IDs (Roles: admin)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        body body RestoreBulkRequest true "Bulk restore request"
// @Success      200 {object} common.SuccessResponse "Products restored successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body"
// @Failure      500 {object} common.ErrorResponse "Failed to restore products"
// @x-roles      ["admin"]
// @Router       /products/trash/restore-bulk [post]
func (h *PrdHandler) RestoreProductsBulkHandler(ctx fiber.Ctx) error {
	var req RestoreBulkRequest
	if err := ctx.Bind().Body(&req); err != nil {
		h.log.Warn("Bulk restore request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	err := h.prdService.RestoreProductsBulk(ctx.RequestCtx(), req)
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

// GetStockHistoryHandler gets stock history for a product
// @Summary      Get stock history for a product
// @Description  Get stock history for a specific product by its ID with pagination (Roles: admin, manager)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id    path  string true  "Product ID" Format(uuid)
// @Param        page  query int    false "Page number"
// @Param        limit query int    false "Limit"
// @Success      200 {object} common.SuccessResponse{data=PagedStockHistoryResponse} "Stock history retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format or query parameters"
// @Failure      404 {object} common.ErrorResponse "Product not found"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve stock history"
// @x-roles      ["admin", "manager"]
// @Router       /products/{id}/stock-history [get]
func (h *PrdHandler) GetStockHistoryHandler(ctx fiber.Ctx) error {
	productID, err := fiber.Convert(ctx.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warn("Invalid product ID format", "error", err, "id", ctx.Params("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid product ID format"})
	}

	var req ListStockHistoryRequest
	if err := ctx.Bind().Query(&req); err != nil {
		h.log.Warn("Stock history request validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	historyResponse, err := h.prdService.GetStockHistory(ctx.RequestCtx(), productID, req)
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
