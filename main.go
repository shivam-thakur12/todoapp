package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Initialize repository and service
	repo := NewTodoRepository()
	service := NewTodoService(repo)
	handler := &TodoHandler{Service: service}

	// Setup routes with initialized handler
	r := setupRoutes(handler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
