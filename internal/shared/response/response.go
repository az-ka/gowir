package response

import (
	"encoding/json"
	"net/http"
)

// BaseResponse adalah format standar semua API response
type BaseResponse struct {
	Status  string      `json:"status"`            // "success" atau "error"
	Message string      `json:"message"`           // Pesan untuk user
	Data    any `json:"data,omitempty"`    // Untuk response sukses
	Errors  any `json:"errors,omitempty"`  // Untuk detail error (misal validasi)
}

// ErrorDetail mendefinisikan struktur error per field (misal validasi)
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// JSON membungkus response sukses
func JSON(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	status := "success"
	if code >= 400 {
		status = "error"
	}

	json.NewEncoder(w).Encode(BaseResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Error membungkus response error tanpa detail
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, message, nil)
}

// ValidationError membungkus response error dengan detail validasi field (Status 422)
func ValidationError(w http.ResponseWriter, message string, errors []ErrorDetail) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity) // 422

	json.NewEncoder(w).Encode(BaseResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}
