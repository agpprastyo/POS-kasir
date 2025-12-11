package common

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"POS-kasir/pkg/logger"
	validatorpkg "POS-kasir/pkg/validator"
)

// ValidateAndRespond melakukan validasi terhadap `req` menggunakan validator `v`.
// Jika validasi gagal, fungsi ini menulis response JSON ke client (400) dan
// mengembalikan (true, errFromFiber). Jika validasi sukses, mengembalikan (false, nil).
//
// - c: fiber context
// - v: instance validator dari pkg/validator (implements Validate(interface{}) error)
// - log: logger untuk mencatat jika penulisan response gagal
// - req: pointer ke struct request yang ingin divalidasi (mis. &req)
//
// Contoh pemakaian di handler:
//
//	if done, err := common.ValidateAndRespond(c, h.Validator, h.Log, &req); done {
//	    return err
//	}
func ValidateAndRespond(c *fiber.Ctx, v validatorpkg.Validator, log logger.ILogger, req interface{}) (bool, error) {
	if err := v.Validate(req); err != nil {
		// coba cast ke structured ValidationErrors dari pkg/validator
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

		// fallback: non-structured validation error
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
