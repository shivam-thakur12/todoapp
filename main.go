package main

import (
	"TODO/pkg/config"
	"TODO/pkg/handlers"
	"TODO/pkg/redis"
	"TODO/pkg/repo"
	"TODO/pkg/routes"
	"TODO/pkg/server"
	"TODO/pkg/service"

	"log"
	"net/http"
)

func main() {
	// Load the configuration
	configg := config.InitConfig()

	// Initialize the database
	server.InitDB(configg) // Pass the config to the initDB function
	defer server.DB.Close()

	// Initialize repository and service
	repo := repo.NewTodoRepository()
	// Initialize Redis client
	client := server.NewRedisClient(configg)
	cache := redis.NewRedisCache(client)

	faktoryClient := server.InitFaktory(configg)

	service := service.NewTodoService(repo, cache, faktoryClient)
	handler := &handlers.TodoHandler{Service: service}

	// Setup routes with initialized handler
	r := routes.SetupRoutes(handler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
