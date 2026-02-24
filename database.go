package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB creates a connection pool to PostgreSQL using pgx
func ConnectDB(ctx context.Context, dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to create database connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database did not respond to ping: %v", err)
	}

	log.Info("successfully connected to PostgreSQL")
	return pool
}