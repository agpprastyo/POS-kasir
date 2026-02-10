package settings

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	service ISettingsService
}

func NewSettingsHandler(service ISettingsService) *SettingsHandler {
	return &SettingsHandler{
		service: service,
	}
}

// GetBrandingHandler godoc
// @Summary Get branding settings
// @Description Get branding settings like app name, logo, footer text
// @Tags Settings
// @Produce json
// @Success 200 {object} dto.BrandingSettingsResponse
// @Router /settings/branding [get]
func (h *SettingsHandler) GetBrandingHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	resp, err := h.service.GetBranding(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to fetch branding settings",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Branding settings fetched successfully",
		Data:    resp,
	})
}

// UpdateBrandingHandler godoc
// @Summary Update branding settings
// @Description Update branding settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body dto.UpdateBrandingRequest true "Update Branding Request"
// @Success 200 {object} dto.BrandingSettingsResponse
// @Router /settings/branding [put]
func (h *SettingsHandler) UpdateBrandingHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.UpdateBrandingRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	resp, err := h.service.UpdateBranding(ctx, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to update branding settings",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Branding settings updated successfully",
		Data:    resp,
	})
}

// UpdateLogoHandler godoc
// @Summary Update app logo
// @Description Upload and update app logo
// @Tags Settings
// @Accept multipart/form-data
// @Produce json
// @Param logo formData file true "Logo image file"
// @Success 200 {object} map[string]string
// @Router /settings/branding/logo [post]
func (h *SettingsHandler) UpdateLogoHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Logo file is required",
			Error:   err.Error(),
		})
	}

	// Validate file size (e.g. max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "File size exceeds 5MB limit",
		})
	}

	// Validate content type
	src, err := file.Open()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to open file",
			Error:   err.Error(),
		})
	}
	defer src.Close()

	// Read file content
	data, err := io.ReadAll(src)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to read file content",
			Error:   err.Error(),
		})
	}

	contentType := file.Header.Get("Content-Type")

	url, err := h.service.UpdateLogo(ctx, data, file.Filename, contentType)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to update logo",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Logo updated successfully",
		Data:    map[string]string{"url": url},
	})
}

// GetPrinterSettingsHandler godoc
// @Summary Get printer settings
// @Description Get printer settings like connection string, paper width
// @Tags Settings
// @Produce json
// @Success 200 {object} dto.PrinterSettingsResponse
// @Router /settings/printer [get]
func (h *SettingsHandler) GetPrinterSettingsHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	resp, err := h.service.GetPrinterSettings(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to fetch printer settings",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Printer settings fetched successfully",
		Data:    resp,
	})
}

// UpdatePrinterSettingsHandler godoc
// @Summary Update printer settings
// @Description Update printer settings
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body dto.UpdatePrinterSettingsRequest true "Update Printer Settings Request"
// @Success 200 {object} dto.PrinterSettingsResponse
// @Router /settings/printer [put]
func (h *SettingsHandler) UpdatePrinterSettingsHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.UpdatePrinterSettingsRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	resp, err := h.service.UpdatePrinterSettings(ctx, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to update printer settings",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Printer settings updated successfully",
		Data:    resp,
	})
}
