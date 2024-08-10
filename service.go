package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type Todo struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type TodoService interface {
	CreateTodoService(body io.Reader) (*Todo, error)
	GetTodosService() ([]Todo, error)
	UpdateTodoService(id int, body io.Reader) (*Todo, error)
	DeleteTodoService(id int, body io.Reader) error
}

type todoService struct {
	Repo TodoRepository
}

func NewTodoService(repo TodoRepository) TodoService {
	return &todoService{
		Repo: repo,
	}
}

// Service function for creating a todo
func (s *todoService) CreateTodoService(body io.Reader) (*Todo, error) {
	var todo Todo
	if err := json.NewDecoder(body).Decode(&todo); err != nil {
		return nil, err
	}
	err := s.Repo.CreateTodoRepo(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// Service function for retrieving todos
func (s *todoService) GetTodosService() ([]Todo, error) {
	return s.Repo.GetTodosRepo()
}

// Service function for updating a todo
func (s *todoService) UpdateTodoService(id int, body io.Reader) (*Todo, error) {
	var updatedTodo Todo
	if err := json.NewDecoder(body).Decode(&updatedTodo); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}
	updatedTodo.ID = id
	rowsAffected, err := s.Repo.UpdateTodoRepo(&updatedTodo)
	if err != nil {
		return nil, fmt.Errorf("database operation error: %w", err)

	}

	// If no rows were affected, it means the Todo with the given ID was not found
	if rowsAffected == 0 {
		return nil, fmt.Errorf("Todo not found")
	}
	// Check if the ID is provided
	if updatedTodo.ID == 0 {
		return nil, fmt.Errorf("Invalid ID")
	}
	return &updatedTodo, nil
}
func (s *todoService) DeleteTodoService(id int, body io.Reader) error {

	rowsAffected, err := s.Repo.DeleteTodoRepo(id)
	if err != nil {
		return fmt.Errorf("database operation error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Todo not found")
	}

	return err
}
