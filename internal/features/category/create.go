package category

import (
	"encoding/json"
	"errors"
	"gowir/internal/db"
	"gowir/internal/shared/response"
	"gowir/internal/shared/util"
	"gowir/internal/shared/validator"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

	// Normalisasi: Hapus spasi di awal dan akhir
	req.Name = strings.TrimSpace(req.Name)
	if req.Description != nil {
		trimmedDesc := strings.TrimSpace(*req.Description)
		req.Description = &trimmedDesc
	}

	// 2. Validate Fields - 422 Unprocessable Entity
	if err := h.validator.Struct(req); err != nil {
		errors := validator.ParseValidationErrors(err)
		response.ValidationError(w, "Beberapa field tidak valid. Silakan periksa kembali input Anda.", errors)
		return
	}

	// 3. Generate ID
	id := util.MustNewUUID()

	// 4. Generate Slug
	categorySlug := util.GenerateSlug(req.Name)

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 409 Conflict - Duplicate slug (SQLSTATE 23505: unique_violation)
			if pgErr.Code == "23505" {
				log.Error("category creation failed: duplicate name", "err", err)
				response.Error(w, 409, "Nama kategori ini sudah digunakan. Silakan gunakan nama lain.")
				return
			}
			// 422 Unprocessable Entity - Invalid parent_id (SQLSTATE 23503: foreign_key_violation)
			if pgErr.Code == "23503" {
				log.Error("category creation failed: invalid parent_id", "parent_id", req.ParentID)
				response.Error(w, 422, "Kategori induk yang Anda pilih tidak tersedia.")
				return
			}
		}
		// 500 Internal Server Error
		log.Error("category creation failed: database error", "err", err)
		response.Error(w, 500, "Maaf, saat ini kami sedang mengalami kendala teknis. Silakan coba beberapa saat lagi.")
		return
	}

	// 7. Success - 201 Created
	response.JSON(w, 201, "Kategori baru berhasil ditambahkan.", category)
}
