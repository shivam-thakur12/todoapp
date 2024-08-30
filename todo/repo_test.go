package todo

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodoRepo(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Expected behavior for the query
	mock.ExpectQuery("INSERT INTO todos").
		WithArgs("Test Todo", "pending").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create a new repository and a Todo object
	todoRepo := NewTodoRepository()
	todo := &Todo{Title: "Test Todo", Status: "pending"}

	// Call the repository function
	err = todoRepo.CreateTodoRepo(todo)

	// Assert that no error occurred and the ID was set correctly
	assert.NoError(t, err)
	assert.Equal(t, 1, todo.ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestCreateTodoRepo_ScanError(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Simulate an error during the Scan() call
	mock.ExpectQuery("INSERT INTO todos").
		WithArgs("Test Todo", "pending").
		WillReturnError(assert.AnError)

	// Create a new repository and a Todo object
	todoRepo := NewTodoRepository()
	todo := &Todo{Title: "Test Todo", Status: "pending"}

	// Call the repository function
	err = todoRepo.CreateTodoRepo(todo)

	// Assert that an error occurred and the ID was not set
	assert.Error(t, err)
	assert.Equal(t, 0, todo.ID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetTodosRepo(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Expected behavior for the query
	rows := sqlmock.NewRows([]string{"id", "title", "status"}).
		AddRow(1, "Test Todo", "pending")

	mock.ExpectQuery("SELECT id, title, status FROM todos WHERE deleted_at IS NULL").
		WillReturnRows(rows)

	// Create a new repository
	todoRepo := NewTodoRepository()

	// Call the repository function
	todos, err := todoRepo.GetTodosRepo()

	// Assert that no error occurred and the results match
	assert.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, 1, todos[0].ID)
	assert.Equal(t, "Test Todo", todos[0].Title)
	assert.Equal(t, "pending", todos[0].Status)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetTodosRepo_QueryError(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Simulate an error during the Query() call
	mock.ExpectQuery("SELECT id, title, status FROM todos WHERE deleted_at IS NULL").
		WillReturnError(assert.AnError)

	// Create a new repository
	todoRepo := NewTodoRepository()

	// Call the repository function
	todos, err := todoRepo.GetTodosRepo()

	// Assert that an error occurred and no todos were returned
	assert.Error(t, err)
	assert.Nil(t, todos)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTodoRepo(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Expected behavior for the query
	mock.ExpectExec("UPDATE todos SET title=\\$1, status=\\$2 WHERE id=\\$3 AND deleted_at IS NULL").
		WithArgs("Updated Todo", "completed", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new repository and a Todo object
	todoRepo := NewTodoRepository()
	todo := &Todo{ID: 1, Title: "Updated Todo", Status: "completed"}

	// Call the repository function
	rowsAffected, err := todoRepo.UpdateTodoRepo(todo)

	// Assert that no error occurred and one row was affected
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateTodoRepo_DBExecError(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Simulate an error when the Exec method is called
	mock.ExpectExec("UPDATE todos SET title=\\$1, status=\\$2 WHERE id=\\$3 AND deleted_at IS NULL").
		WithArgs("Updated Todo", "completed", 1).
		WillReturnError(errors.New("database execution error"))

	// Create a new repository and a Todo object
	todoRepo := NewTodoRepository()
	todo := &Todo{ID: 1, Title: "Updated Todo", Status: "completed"}

	// Call the repository function
	rowsAffected, err := todoRepo.UpdateTodoRepo(todo)

	// Assert that an error occurred and no rows were affected
	assert.Error(t, err)
	assert.Equal(t, int64(0), rowsAffected)
	assert.Equal(t, "database execution error", err.Error())

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTodoRepo_RowsAffectedError(t *testing.T) {
	// Create a new mock database connection and mock object
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mocking the DB variable from the server package
	DB = db

	// Simulate an error in RowsAffected
	mock.ExpectExec("UPDATE todos SET title=\\$1, status=\\$2 WHERE id=\\$3 AND deleted_at IS NULL").
		WithArgs("Updated Todo", "completed", 1).
		WillReturnResult(sqlmock.NewErrorResult(assert.AnError)) // Simulate an error here

	// Create a new repository and a Todo object
	todoRepo := NewTodoRepository()
	todo := &Todo{ID: 1, Title: "Updated Todo", Status: "completed"}

	// Call the repository function
	rowsAffected, err := todoRepo.UpdateTodoRepo(todo)

	// Assert that an error occurred and no rows were affected
	assert.Error(t, err)                    // Expecting an error here
	assert.Equal(t, int64(0), rowsAffected) // No rows should be affected

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
