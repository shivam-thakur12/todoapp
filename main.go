package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	r := setupRoutes()
	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))

}
