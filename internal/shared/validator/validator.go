package validator

import (
	"fmt"
	"gowir/internal/shared/response"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ParseValidationErrors menerjemahkan error dari library validator ke format JSON API kita
func ParseValidationErrors(err error) []response.ErrorDetail {
	var details []response.ErrorDetail

	// Jika error adalah bertipe ValidationErrors dari library
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			field := strings.ToLower(err.Field())
			message := ""

			// Terjemahkan tag validator ke Bahasa Indonesia
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("Field %s wajib diisi", field)
			case "min":
				message = fmt.Sprintf("Field %s minimal %s karakter", field, err.Param())
			case "max":
				message = fmt.Sprintf("Field %s maksimal %s karakter", field, err.Param())
			case "email":
				message = fmt.Sprintf("Field %s harus berupa alamat email yang valid", field)
			case "uuid":
				message = fmt.Sprintf("Field %s harus berupa UUID yang valid", field)
			default:
				message = fmt.Sprintf("Field %s tidak valid", field)
			}

			details = append(details, response.ErrorDetail{
				Field:   field,
				Message: message,
			})
		}
	}

	return details
}
