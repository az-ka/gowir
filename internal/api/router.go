package api

import (
	"net/http"

	"gowir/internal/db"
	"gowir/internal/features/category"
	"gowir/internal/shared/response"
	"gowir/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

// NewRouter bertugas merakit dan mendaftarkan semua endpoint aplikasi
func NewRouter(queries *db.Queries, validate *validator.Validate) *chi.Mux {
	r := chi.NewRouter()

	// 1. Global Middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.MiddlewareLogging)

	// 2. Health Check (Root)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, 200, "API e-commerce berjalan dengan baik", nil)
	})

	// 3. Inisialisasi semua Handler
	categoryHandler := category.NewHandler(queries, validate)

	// 4. Daftarkan Routes berdasarkan Group
	r.Route("/api/v1", func(r chi.Router) {
		// Group: ADMIN
		r.Route("/admin", func(r chi.Router) {
			// Pasang middleware pelindung untuk semua endpoint admin
			r.Use(middleware.RequireAuth)
			r.Use(middleware.RequireAdmin)

			categoryHandler.RegisterRoutes(r)
		})
	})

	return r
}
