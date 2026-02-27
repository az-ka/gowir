package category

import (
	"gowir/internal/db"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	queries   *db.Queries
	validator *validator.Validate
}

func NewHandler(queries *db.Queries, validate *validator.Validate) *Handler {
	return &Handler{
		queries:   queries,
		validator: validate,
	}
}
