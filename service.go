package main

import (
	"encoding/json"
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
	UpdateTodoService(body io.Reader) (*Todo, error)
	DeleteTodoService(body io.Reader) (int, error)
	CreateTodoRepo(todo *Todo) error
	GetTodosRepo() ([]Todo, error)
	UpdateTodoRepo(todo *Todo) (int64, error)
	DeleteTodoRepo(id int) (int64, error)
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
	return &todo, nil
}

// Service function for retrieving todos
func (s *todoService) GetTodosService() ([]Todo, error) {
	return s.Repo.GetTodosRepo()
}

// Service function for updating a todo
func (s *todoService) UpdateTodoService(body io.Reader) (*Todo, error) {
	var updatedTodo Todo
	if err := json.NewDecoder(body).Decode(&updatedTodo); err != nil {
		return nil, err
	}

	return &updatedTodo, nil
}

// Service function for deleting a todo
func (s *todoService) DeleteTodoService(body io.Reader) (int, error) {
	var todo Todo
	if err := json.NewDecoder(body).Decode(&todo); err != nil {
		return 0, err
	}
	return todo.ID, nil
}

// Methods to call repo functions
func (s *todoService) CreateTodoRepo(todo *Todo) error {
	return s.Repo.CreateTodoRepo(todo)
}

func (s *todoService) GetTodosRepo() ([]Todo, error) {
	return s.Repo.GetTodosRepo()
}

func (s *todoService) UpdateTodoRepo(todo *Todo) (int64, error) {
	return s.Repo.UpdateTodoRepo(todo)
}

func (s *todoService) DeleteTodoRepo(id int) (int64, error) {
	return s.Repo.DeleteTodoRepo(id)
}
