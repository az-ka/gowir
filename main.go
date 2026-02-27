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

	pool, err := ConnectDB(ctx, dbURL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
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
		fmt.Fprintln(w, "API e-commerce berjalan dengan baik")
	})

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Info("server is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Error("server encountered a fatal error", "error", err)
	case sig := <-quit:
		log.Info("received termination signal", "signal", sig)
	}

	log.Info("shutting down server gracefully...")

	ctxShutdown, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server stopped safely")
}
