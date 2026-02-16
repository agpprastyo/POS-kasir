package common

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"POS-kasir/pkg/logger"
	validatorpkg "POS-kasir/pkg/validator"
)

func ValidateAndRespond(c fiber.Ctx, v validatorpkg.Validator, log logger.ILogger, req interface{}) (bool, error) {
	if err := v.Validate(req); err != nil {

		var ve *validatorpkg.ValidationErrors
		if errors.As(err, &ve) {
			resp := ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			}
			errResp := c.Status(fiber.StatusBadRequest).JSON(resp)
			if errResp != nil {
				// log jika gagal menulis response (mis. connection closed)
				log.Errorf("ValidateAndRespond: failed to write validation response: %v", errResp)
			}
			return true, errResp
		}

		resp := ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		}
		errResp := c.Status(fiber.StatusBadRequest).JSON(resp)
		if errResp != nil {
			log.Errorf("ValidateAndRespond: failed to write validation response: %v", errResp)
		}
		return true, errResp
	}
	return false, nil
}
