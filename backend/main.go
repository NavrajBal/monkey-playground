package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"monkey-playground-backend/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Routes
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API routes
	mux.HandleFunc("/api/tokenize", api.TokenizeHandler)
	mux.HandleFunc("/api/parse", api.ParseHandler)
	mux.HandleFunc("/api/compile", api.CompileHandler)
	mux.HandleFunc("/api/execute", api.ExecuteHandler)
	mux.HandleFunc("/api/repl", api.ReplHandler)

	// Apply CORS middleware
	handler := corsHandler(mux)

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /health")
	fmt.Println("  POST /api/tokenize")
	fmt.Println("  POST /api/parse")
	fmt.Println("  POST /api/compile")
	fmt.Println("  POST /api/execute")
	fmt.Println("  POST /api/repl")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
