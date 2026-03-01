package request

import (
	"encoding/json"
	"errors"
	"gowir/internal/shared/response"
	"io"
	"net/http"
)

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		response.Error(w, 400, "Format data yang dikirim tidak valid atau terdapat field yang tidak dikenali.")
		return false
	}

	// Pastikan tidak ada data tambahan setelah JSON selesai dibaca
	err := decoder.Decode(&struct{}{})
	if err != nil && !errors.Is(err, io.EOF) {
		response.Error(w, 400, "Format data yang dikirim tidak valid: terdapat struktur JSON ganda.")
		return false
	}

	return true
}
