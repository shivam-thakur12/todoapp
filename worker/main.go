package main

import (
	"TODO/pkg/config"
	"TODO/pkg/server"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the configuration
	configg := config.InitConfig()

	// Initialize the database connection
	server.InitDB(configg)
	defer server.DB.Close()

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
