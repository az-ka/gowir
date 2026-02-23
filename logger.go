package main

import (
	"os"

	"github.com/charmbracelet/log"
)

func init() {
	if os.Getenv("LOG_LEVEL") == "debug" {
		log.SetLevel(log.DebugLevel)
	}
}
