package request

import (
	"encoding/json"
	"gowir/internal/shared/response"
	"net/http"
)

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, 400, "Format data yang dikirim tidak valid. Harap periksa format JSON Anda.")
		return false
	}
	return true
}
