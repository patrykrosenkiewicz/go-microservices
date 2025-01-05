package main

import (
	"fmt"
	"log"

	"github.com/patrykrosenkiewicz/go-microservices/http"
)

func main() {

	port := "8080"
	// Create a new server
	s := http.NewServer(port)

	// Start the HTTP server
	fmt.Println("Server is running on port %s...", port)

	// Start the HTTP server.
	if err := s.Open(); err != nil {
		// Handle the error properly, since main() doesn't return any value.
		log.Fatalf("Error starting server: %v", err)
	}
}
