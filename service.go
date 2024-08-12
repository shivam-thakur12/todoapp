package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
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
	Repo  TodoRepository
	Cache RedisCache
}

func NewTodoService(repo TodoRepository, cache RedisCache) TodoService {
	return &todoService{
		Repo:  repo,
		Cache: cache,
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
	s.Cache.Del("todos_list") // Invalidate cache
	return &todo, nil
}

// Service function for retrieving todos
func (s *todoService) GetTodosService() ([]Todo, error) {
	// Define a cache key for todos
	cacheKey := "todos_list"

	// Attempt to retrieve todos from Redis cache
	cachedTodos, err := s.Cache.Get(cacheKey)
	if err == nil && cachedTodos != "" {
		var todos []Todo
		// Deserialize JSON from cache into todos slice
		if err := json.Unmarshal([]byte(cachedTodos), &todos); err == nil {
			// fmt.Println("cache is working")
			// Return the todos from cache
			return todos, nil
		}
	}

	// If cache miss or error, retrieve todos from database
	todos, err := s.Repo.GetTodosRepo()
	if err != nil {
		return nil, fmt.Errorf("database operation error: %w", err)
	}

	// Serialize the todos into JSON and cache it in Redis for 10 minutes
	todosJSON, err := json.Marshal(todos)
	if err == nil {
		s.Cache.Set(cacheKey, todosJSON, 10*time.Minute)
		// fmt.Println("cache set is working")
	}
	return todos, nil
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
	s.Cache.Del("todos_list") // Invalidate cache
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
	s.Cache.Del("todos_list") // Invalidate cache
	return err
}
