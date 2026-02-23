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
	if err := godotenv.Load(); err != nil {
		log.Debug("failed to load .env", "err", err)
	}

	dbURL = os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not found")
	}

	var ok bool
	port, ok = os.LookupEnv("APP_PORT")
	if !ok || port == "" {
		port = ":4000"
	}
	if port[0] != ':' {
		port = ":" + port
	}
}
