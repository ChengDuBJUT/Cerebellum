package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cerebellum/internal/brain"
	"cerebellum/internal/config"
	"cerebellum/internal/llm"
	"cerebellum/internal/server"
	"cerebellum/internal/store"
)

func main() {
	// Load configuration
	cfg, err := config.Load("cerebellum.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize markdown store
	markdownStore, err := store.NewMarkdownStore("brain.md")
	if err != nil {
		log.Fatalf("Failed to create markdown store: %v", err)
	}

	// Initialize LLM client
	llmClient := llm.NewOllama(cfg.Ollama.Host, cfg.Ollama.Model)

	// Initialize file watcher
	watcher, err := brain.NewWatcher("brain.md", markdownStore)
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	go watcher.Start()

	addr := cfg.GetServerAddr()
	log.Printf("Cerebellum server starting on %s", addr)

	// Create server
	httpServer := server.NewServer(cfg, markdownStore, llmClient)

	// Use custom mux
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Got /api/status request: %s %s", r.Method, r.URL.Path)
		httpServer.HandleAPIStatus(w, r)
	})
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Got /api/tasks request: %s %s", r.Method, r.URL.Path)
		httpServer.HandleAPITasks(w, r)
	})
	mux.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Got /api/report request: %s %s", r.Method, r.URL.Path)
		httpServer.HandleAPIReport(w, r)
	})
	mux.HandleFunc("/health", httpServer.HandleHealth)
	mux.HandleFunc("/tasks", httpServer.HandleTasks)
	mux.HandleFunc("/reload", httpServer.HandleReload)
	mux.HandleFunc("/api/task/", httpServer.HandleAPITaskDelete)
	mux.HandleFunc("/chat", httpServer.HandleChat)
	mux.HandleFunc("/execute", httpServer.HandleExecute)
	mux.HandleFunc("/api/chat", httpServer.HandleChat)
	mux.HandleFunc("/api/execute", httpServer.HandleExecute)

	log.Printf("DEBUG: Mux handlers registered, addr=%s", addr)

	// Start task executor
	go httpServer.StartTaskExecutor()

	// Start HTTP server
	go func() {
		log.Printf("DEBUG: Starting ListenAndServe on %s", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Cerebellum server started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	watcher.Stop()
	log.Println("Stopped")
}
