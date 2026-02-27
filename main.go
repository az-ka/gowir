package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"gowir/internal/db"
	"gowir/internal/features/category"
	"gowir/middleware"
)

func main() {
	Config()

	ctx := context.Background()
	pool := ConnectDB(ctx, dbURL)
	defer pool.Close()

	// Init SQLC queries to be passed to handlers
	queries := db.New(pool)
	_ = queries

	validate := validator.New()

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.MiddlewareLogging)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/admin", func(r chi.Router) {
			categoryHandler := category.NewHandler(queries, validate)
			categoryHandler.RegisterRoutes(r)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API e-commerce berjalan dengan baik") // User facing message
	})

	log.Info("server is running", "port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
