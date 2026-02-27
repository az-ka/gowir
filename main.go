package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := ConnectDB(ctx, dbURL)
	defer pool.Close()

	// Init SQLC queries to be passed to handlers
	queries := db.New(pool)
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

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("server is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("gagal menjalankan server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("mematikan server secara perlahan...")

	ctxShutdown, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("server terpaksa dimatikan", "error", err)
	}

	log.Info("server berhasil berhenti dengan aman")
}
