package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

func init() {
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
