package category

import (
	"encoding/json"
	"gowir/internal/db"
	"gowir/internal/shared/response"
	"gowir/internal/shared/validator"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type CreateCategoryReq struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description *string    `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryReq

	// 1. Decode JSON - 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "Format data yang dikirim tidak valid. Harap periksa format JSON Anda.")
		return
	}

	// 2. Validate Fields - 422 Unprocessable Entity
	if err := h.validator.Struct(req); err != nil {
		errors := validator.ParseValidationErrors(err)
		response.ValidationError(w, "Beberapa field tidak valid. Silakan periksa kembali input Anda.", errors)
		return
	}

	// 3. Generate ID - 500 Internal Server Error
	id, err := uuid.NewV7()
	if err != nil {
		response.Error(w, 500, "Gagal membuat identitas kategori yang unik. Silakan coba lagi.")
		return
	}

	// 4. Generate Slug
	categorySlug := generateCustomSlug(req.Name)

	// 5. Save to Database
	category, err := h.queries.CreateCategory(r.Context(), db.CreateCategoryParams{
		ID:          id,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        categorySlug,
		Description: req.Description,
	})

	// 6. Handle Database Errors
	if err != nil {
		// 409 Conflict - Duplicate slug
		if strings.Contains(err.Error(), "unique constraint") {
			response.Error(w, 409, "Kategori dengan nama atau slug tersebut sudah ada.")
			return
		}
		// 500 Internal Server Error
		response.Error(w, 500, "Terjadi kesalahan saat menyimpan kategori ke database. Silakan coba lagi nanti.")
		return
	}

	// 7. Success - 201 Created
	response.JSON(w, 201, "Kategori baru berhasil ditambahkan.", category)
}

func generateCustomSlug(name string) string {
	// Ganti symbol "&" menjadi "dan" agar lebih user-friendly
	formatted := strings.ReplaceAll(name, "&", "dan")
	return slug.Make(formatted)
}
