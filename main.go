package main

import (
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the database
	initDB()
	defer db.Close()
	setupRoutes()

}
