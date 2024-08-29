package main

import (
	"TODO/todo/server"
	"context"
	"fmt"
	"log"

	worker "github.com/contribsys/faktory_worker_go"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func deleteTodoWorker(ctx context.Context, args ...interface{}) error {
	help := worker.HelperFor(ctx)
	log.Printf("Working on job %s\n", help.Jid())

	// Extract the ID from job arguments
	idMap := args[0].(map[string]interface{})
	id := int(idMap["id"].(float64))

	// Perform the delete operation
	res, err := server.DB.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Todo with ID %d not found or already deleted\n", id)
		// Stop retries by returning nil
		return nil
	}

	log.Printf("Todo with ID %d marked as deleted\n", id)
	return nil

}
