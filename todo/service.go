package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/contribsys/faktory/client"
	"github.com/google/uuid"
)

type TodoService interface {
	CreateTodoService(body io.Reader) (*Todo, error)
	GetTodosService() ([]Todo, error)
	UpdateTodoService(id int, body io.Reader) (*Todo, error)
	DeleteTodoService(id int, body io.Reader) error
}

type todoService struct {
	Repo          TodoRepository
	Cache         RedisCache
	FaktoryClient *client.Client
}

func NewTodoService(repo TodoRepository, cache RedisCache, faktoryClient *client.Client) TodoService {
	return &todoService{
		Repo:          repo,
		Cache:         cache,
		FaktoryClient: faktoryClient,
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
	err = s.Cache.Del(CacheKeyTodosList) // Invalidate cache
	if err != nil {
		return nil, fmt.Errorf("error invalidating cache %v", err)
	}
	return &todo, nil
}

// Service function for retrieving todos
func (s *todoService) GetTodosService() ([]Todo, error) {

	// Attempt to retrieve todos from cache
	todos, err := s.retrieveTodosFromCache(CacheKeyTodosList)
	if err == nil && todos != nil {
		// Return the todos from cache
		return todos, nil
	}

	// If cache miss or error, retrieve todos from database
	todos, err = s.Repo.GetTodosRepo()
	if err != nil {
		return nil, fmt.Errorf("database operation error: %w", err)
	}

	// Cache the todos for future requests
	err = s.cacheTodos(todos, CacheKeyTodosList)
	if err != nil {
		// Optionally log the error but continue; you still return the data from DB
		fmt.Printf("Warning: failed to cache todos: %v\n", err)
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
	// // Check if the ID is provided
	// if updatedTodo.ID == 0 {
	// 	return nil, fmt.Errorf("Invalid ID")
	// }
	err = s.Cache.Del(CacheKeyTodosList) // Invalidate cache
	if err != nil {
		return nil, fmt.Errorf("error invalidating cache %v", err)
	}
	return &updatedTodo, nil
}
func (s *todoService) DeleteTodoService(id int, body io.Reader) error {

	// Generate a unique job ID
	jid := uuid.NewString()

	job := &client.Job{
		Queue: "default",
		Type:  "delete_todo",
		Args:  []interface{}{map[string]interface{}{"id": id}},
		Jid:   jid, // Set the unique JID
	}

	// Push job to Faktory
	err := s.FaktoryClient.Push(job)
	if err != nil {
		return fmt.Errorf("failed to push job to Faktory: %w", err)
	}

	// Optionally invalidate cache if needed
	err = s.Cache.Del(CacheKeyTodosList)
	if err != nil {
		return fmt.Errorf("error invalidating cache %v", err)
	}
	return nil
}

// Service function to cache todos
func (s *todoService) cacheTodos(todos []Todo, cacheKey string) error {
	// Serialize the todos into JSON
	todosJSON, err := json.Marshal(todos)
	if err != nil {
		return fmt.Errorf("failed to marshal todos to JSON: %w", err)
	}

	// Set the cached value in Redis
	err = s.Cache.Set(cacheKey, todosJSON, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to set cache in Redis: %w", err)
	}
	// fmt.Println("cached")
	return nil
}

// Service function to retrieve todos from cache
func (s *todoService) retrieveTodosFromCache(cacheKey string) ([]Todo, error) {
	// Attempt to retrieve todos from Redis cache
	cachedTodos, err := s.Cache.Get(cacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache from Redis: %w", err)
	}
	// if cachedTodos == "" {
	// 	return nil, nil // Cache miss
	// }

	var todos []Todo
	// Deserialize JSON from cache into todos slice
	err = json.Unmarshal([]byte(cachedTodos), &todos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached todos: %w", err)
	}
	// fmt.Println("retreived from cache")
	return todos, nil
}
