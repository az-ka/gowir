package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"

	"gowir/internal/api"
	"gowir/internal/db"
)

func main() {
	Config()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := ConnectDB(ctx, dbURL)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return
	}
	defer pool.Close()

	// Init SQLC queries to be passed to handlers
	queries := db.New(pool)
	validate := validator.New()

	r := api.NewRouter(queries, validate)

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
