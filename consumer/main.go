package main

import (
	"TODO/consumer/worker"
	"TODO/todo"
	"TODO/todo/config"
	"log"
	"time"

	workerr "github.com/contribsys/faktory_worker_go"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the configuration
	configg := config.InitConfig()

	// Initialize the database connection
	todo.InitDB(configg)
	defer todo.DB.Close()

	// Initialize Faktory worker manager
	mgr := workerr.NewManager()

	// Register the delete job
	mgr.Register("delete_todo", worker.DeleteTodoWorker)

	// Set concurrency and shutdown timeout
	mgr.Concurrency = 1
	mgr.ShutdownTimeout = 25 * time.Second
	// Process jobs from these queues
	mgr.ProcessStrictPriorityQueues("default")

	// Start processing jobs
	err := mgr.Run()
	if err != nil {
		log.Fatal(err)
	}
}
