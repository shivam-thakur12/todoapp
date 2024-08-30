package main

import (
	"TODO/cmd/api"
	"TODO/todo"
	"TODO/todo/config"
	"log"
	"net/http"
)

func main() {
	// Load the configuration
	configg := config.InitConfig()

	// Initialize the database
	todo.InitDB(configg) // Pass the config to the initDB function
	todo.RunMigrations(configg)
	defer todo.DB.Close()

	// Initialize repository and service
	repo := todo.NewTodoRepository()
	// Initialize Redis client
	client := todo.NewRedisClient(configg)
	cache := todo.NewRedisCache(client)

	faktoryClient := todo.InitFaktory(configg)

	service := todo.NewTodoService(repo, cache, faktoryClient)
	handler := &api.TodoHandler{Service: service}

	// Setup routes with initialized handler
	r := api.SetupRoutes(handler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
