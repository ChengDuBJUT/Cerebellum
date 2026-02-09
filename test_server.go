package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request: %s %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello from test server on 0.0.0.0:18080")
	})
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got API request: %s %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "API Status OK")
	})
	log.Printf("Starting server on 0.0.0.0:18080")
	if err := http.ListenAndServe("0.0.0.0:18080", mux); err != nil {
		log.Printf("Error: %v", err)
	}
}
