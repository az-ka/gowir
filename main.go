package main

import (
	"fmt"
	"gowir/middleware"
	"net/http"

	"github.com/charmbracelet/log"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":4000", middleware.MiddlewareLogging(mux)))
}
