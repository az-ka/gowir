package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

var (
	dbURL string
	port  string
)

func Config() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL tidak ditemukan di .env")
	}

	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = ":4000"
	}
	if port[0] != ':' {
		port = ":" + port
	}
}
