package main

import (
	"time"

	worker "github.com/contribsys/faktory_worker_go"
)

func main() {
	// Load the configuration
	config := initConfig()

	// Initialize the database connection
	initDB(config)

	// Initialize Faktory worker manager
	mgr := worker.NewManager()

	// Register the delete job
	mgr.Register("delete_todo", deleteTodoWorker)

	// Set concurrency and shutdown timeout
	mgr.Concurrency = 1
	mgr.ShutdownTimeout = 25 * time.Second

	// Process jobs from these queues
	mgr.ProcessStrictPriorityQueues("default")

	// Start processing jobs
	mgr.Run()
}
