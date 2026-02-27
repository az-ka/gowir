package response

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

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

func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, message, nil)
}
