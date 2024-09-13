package worker

import (
	"fmt"
	"testing"
	"time"

	"TODO/todo"
	"TODO/todo/config"

	"github.com/contribsys/faktory/client"
	workerr "github.com/contribsys/faktory_worker_go"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() func() {
	// Initialize the configuration
	configg := config.InitConfig()

	// Initialize the database connection
	todo.InitDB(configg)

	// Initialize Faktory client
	faktoryClient := todo.InitFaktory(configg)

	return func() {
		// Cleanup resources
		todo.DB.Close()
		faktoryClient.Close()
	}
}

func TestDeleteTodoWorker(t *testing.T) {
	// Setup test environment
	teardown := setup()
	defer teardown()

	// Create a Faktory manager and register the worker
	mgr := workerr.NewManager()
	mgr.Register("delete_todo", DeleteTodoWorker)
	mgr.Concurrency = 1
	mgr.ShutdownTimeout = 25 * time.Second

	// Generate a unique job ID
	jid := uuid.NewString()

	// Create a Faktory job
	job := &client.Job{
		Queue: "default",
		Type:  "delete_todo",
		Args:  []interface{}{map[string]interface{}{"id": 1}},
		Jid:   jid,
	}

	// Push the job to Faktory
	err := faktoryPush(mgr, job)
	require.NoError(t, err, "Failed to push job to Faktory")

	// Create a channel to capture any errors from mgr.Run()
	errChan := make(chan error, 1)
	// Run the Faktory worker manager in a goroutine and capture any errors
	go func() {
		errChan <- mgr.Run()
	}()
	// Check for any errors returned by mgr.Run()
	select {
	case err := <-errChan:
		require.NoError(t, err, "Faktory manager encountered an error")
	default:
		// No errors encountered
	}
	// defer mgr.Shutdown() // Ensure the manager is shutdown after the test

	// Wait for a short duration to allow job processing
	time.Sleep(2 * time.Second)

	// Check if the job was processed correctly
	var result int
	err = todo.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id=$1 AND deleted_at IS NOT NULL", 1).Scan(&result)
	require.NoError(t, err, "Failed to query database")
	assert.Equal(t, 1, result, "Expected one row to be marked as deleted")
}

// Helper function to push jobs to Faktory
func faktoryPush(mgr *workerr.Manager, job *client.Job) error {
	if viper.GetBool("faktory_inline") {
		return syntheticPush(mgr, job)
	}
	return realPush(job)
}

func syntheticPush(mgr *workerr.Manager, job *client.Job) error {
	// Inline dispatch using the worker manager
	err := mgr.InlineDispatch(job)
	if err != nil {
		return fmt.Errorf("syntheticPush failed: %w", err)
	}
	return nil
}

func realPush(job *client.Job) error {
	// Real push to Faktory
	client, err := client.Open()
	if err != nil {
		return fmt.Errorf("failed to open Faktory client connection: %w", err)
	}
	err = client.Push(job)
	if err != nil {
		return fmt.Errorf("failed to enqueue Faktory job: %w", err)
	}
	return nil
}
