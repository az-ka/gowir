package category

import (
	"encoding/json"
	"gowir/internal/db"
	"gowir/internal/shared/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type CreateCategoryReq struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description string     `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "Format data yang dikirim tidak valid")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, 400, "Validasi gagal: Silakan periksa kembali input Anda")
		return
	}

	id, _ := uuid.NewV7()
	categorySlug := generateCustomSlug(req.Name)

	category, err := h.queries.CreateCategory(r.Context(), db.CreateCategoryParams{
		ID:          id,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        categorySlug,
		Description: &req.Description,
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			response.Error(w, 409, "Kategori dengan nama tersebut sudah ada")
			return
		}
		response.Error(w, 500, "Terjadi kesalahan pada sistem, silakan coba lagi nanti")
		return
	}

	response.JSON(w, 201, "Kategori berhasil ditambahkan", category)
}

func generateCustomSlug(name string) string {
	formatted := strings.ReplaceAll(name, "&", "dan")
	return slug.Make(formatted)
}
