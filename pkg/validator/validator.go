package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationErrors is an error wrapper that contains per-field messages.
type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}

func (v *ValidationErrors) Error() string {
	parts := make([]string, 0, len(v.Errors))
	for k, msg := range v.Errors {
		parts = append(parts, k+": "+msg)
	}
	return "validation error: " + strings.Join(parts, ", ")
}

// Validator is the public interface used by the app.
type Validator interface {
	Validate(i interface{}) error
}

type DefaultValidator struct {
	validate *validator.Validate
}

// NewValidator creates and configures a validator instance.
func NewValidator() Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}

		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		return name
	})

	// Tambahkan custom validations di sini bila perlu
	// v.RegisterValidation("password", passwordValidationFunc)

	return &DefaultValidator{validate: v}
}

func (v *DefaultValidator) Validate(i interface{}) error {
	if i == nil {
		return nil
	}
	if err := v.validate.Struct(i); err != nil {
		// convert to ValidationErrors when the error type is validator.ValidationErrors
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string, len(ve))
			for _, fe := range ve {
				// fe.Field() will be mapped to JSON tag because of RegisterTagNameFunc
				field := fe.Field()
				// build human friendly message (you can localize this later)
				out[field] = buildErrorMessage(fe)
			}
			return &ValidationErrors{Errors: out}
		}
		// unexpected error from validator
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// buildErrorMessage creates a readable message from a FieldError.
// Customize messages per tag if you need i18n/translation.
func buildErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "length must be at least " + fe.Param()
	case "max":
		return "length must be at most " + fe.Param()
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	case "lte":
		return "must be less than or equal to " + fe.Param()
	case "oneof":
		return "must be one of the allowed values"
	case "uuid":
		return "must be a valid UUID"
	case "url":
		return "must be a valid URL"
	
	default:
		return fe.Error() // fallback to a default message
	}
}
