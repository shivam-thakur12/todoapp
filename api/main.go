package main

import (
	"TODO/api/handlers"
	"TODO/api/routes"
	"TODO/todo/config"
	"TODO/todo/redis"
	"TODO/todo/repo"
	"TODO/todo/server"
	"TODO/todo/service"

	// "TODO/todo/service"

	"log"
	"net/http"
)

func main() {
	// Load the configuration
	configg := config.InitConfig()

	// Initialize the database
	server.InitDB(configg) // Pass the config to the initDB function
	server.RunMigrations(configg)
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
