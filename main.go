package main

import (
	"log"
	"net/http"
)

func main() {
	// Load the configuration
	config := initConfig()

	// Initialize the database
	initDB(config) // Pass the config to the initDB function
	defer db.Close()

	// Initialize repository and service
	repo := NewTodoRepository()
	// Initialize Redis client
	client := NewRedisClient(config)
	cache := NewRedisCache(client)
	service := NewTodoService(repo, cache)
	handler := &TodoHandler{Service: service}

	// Setup routes with initialized handler
	r := setupRoutes(handler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
