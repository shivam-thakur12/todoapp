package main

import (
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
	res, err := db.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Todo not found")
	}
	// s.Cache.Del(CacheKeyTodosList) // Invalidate cache

	log.Printf("Todo with ID %d marked as deleted\n", id)
	return nil
}
