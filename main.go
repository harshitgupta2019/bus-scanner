package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// homeHandler handles the root endpoint
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Welcome to Bus Booking Aggregator API",
		Data:    map[string]string{"version": "1.0.0"},
	}
	sendJSON(w, http.StatusOK, response)
}

// healthHandler handles health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Service is healthy",
	}
	sendJSON(w, http.StatusOK, response)
}

func main() {
	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)

	// Server configuration
	port := ":8080"

	fmt.Printf("🚌 Bus Booking Aggregator starting on port %s\n", port)
	fmt.Printf("📍 http://localhost%s\n", port)

	// Start the server
	log.Fatal(http.ListenAndServe(port, mux))
}
