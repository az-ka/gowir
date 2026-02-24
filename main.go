package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"

	"gowir/internal/db"
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

	r := chi.NewRouter()
	r.Use(middleware.MiddlewareLogging)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API e-commerce berjalan dengan baik") // User facing message
	})

	log.Info("server is running", "port", port)
	log.Fatal(http.ListenAndServe(port, r))
}

