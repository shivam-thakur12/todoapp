<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/service/service.go
>>>>>>> 4fe2e1c (unit tests)
<<<<<<<< HEAD:todo/service.go
package todo
========
package service
>>>>>>>> a07f97f (worker with refactoring code):pkg/service/service.go
<<<<<<< HEAD

import (
	"TODO/pkg/constants"
	"TODO/pkg/redis"
	"TODO/pkg/repo"
=======
========
package todo
>>>>>>>> 4fe2e1c (unit tests):todo/service.go

import (
>>>>>>> 4fe2e1c (unit tests)
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/contribsys/faktory/client"
	"github.com/google/uuid"
)

type TodoService interface {
<<<<<<< HEAD
	CreateTodoService(body io.Reader) (*repo.Todo, error)
	GetTodosService() ([]repo.Todo, error)
	UpdateTodoService(id int, body io.Reader) (*repo.Todo, error)
=======
	CreateTodoService(body io.Reader) (*Todo, error)
	GetTodosService() ([]Todo, error)
	UpdateTodoService(id int, body io.Reader) (*Todo, error)
>>>>>>> 4fe2e1c (unit tests)
	DeleteTodoService(id int, body io.Reader) error
}

type todoService struct {
<<<<<<< HEAD
	Repo          repo.TodoRepository
	Cache         redis.RedisCache
	FaktoryClient *client.Client
}

func NewTodoService(repo repo.TodoRepository, cache redis.RedisCache, faktoryClient *client.Client) TodoService {
=======
	Repo          TodoRepository
	Cache         RedisCache
	FaktoryClient *client.Client
}

func NewTodoService(repo TodoRepository, cache RedisCache, faktoryClient *client.Client) TodoService {
>>>>>>> 4fe2e1c (unit tests)
	return &todoService{
		Repo:          repo,
		Cache:         cache,
		FaktoryClient: faktoryClient,
	}
}

// Service function for creating a todo
<<<<<<< HEAD
func (s *todoService) CreateTodoService(body io.Reader) (*repo.Todo, error) {
	var todo repo.Todo
=======
func (s *todoService) CreateTodoService(body io.Reader) (*Todo, error) {
	var todo Todo
>>>>>>> 4fe2e1c (unit tests)
	if err := json.NewDecoder(body).Decode(&todo); err != nil {
		return nil, err
	}
	err := s.Repo.CreateTodoRepo(&todo)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
	s.Cache.Del(constants.CacheKeyTodosList) // Invalidate cache
=======
	s.Cache.Del(CacheKeyTodosList) // Invalidate cache
>>>>>>> 4fe2e1c (unit tests)
	return &todo, nil
}

// Service function for retrieving todos
<<<<<<< HEAD
func (s *todoService) GetTodosService() ([]repo.Todo, error) {

	// Attempt to retrieve todos from cache
	todos, err := s.retrieveTodosFromCache(constants.CacheKeyTodosList)
=======
func (s *todoService) GetTodosService() ([]Todo, error) {

	// Attempt to retrieve todos from cache
	todos, err := s.retrieveTodosFromCache(CacheKeyTodosList)
>>>>>>> 4fe2e1c (unit tests)
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
<<<<<<< HEAD
	err = s.cacheTodos(todos, constants.CacheKeyTodosList)
=======
	err = s.cacheTodos(todos, CacheKeyTodosList)
>>>>>>> 4fe2e1c (unit tests)
	if err != nil {
		// Optionally log the error but continue; you still return the data from DB
		fmt.Printf("Warning: failed to cache todos: %v\n", err)
	}

	return todos, nil
}

// Service function for updating a todo
<<<<<<< HEAD
func (s *todoService) UpdateTodoService(id int, body io.Reader) (*repo.Todo, error) {
	var updatedTodo repo.Todo
=======
func (s *todoService) UpdateTodoService(id int, body io.Reader) (*Todo, error) {
	var updatedTodo Todo
>>>>>>> 4fe2e1c (unit tests)
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
<<<<<<< HEAD
<<<<<<<< HEAD:todo/service.go
=======
<<<<<<<< HEAD:todo/service/service.go
<<<<<<<< HEAD:todo/service.go
========
>>>>>>>> 4fe2e1c (unit tests):todo/service.go
>>>>>>> 4fe2e1c (unit tests)
	// // Check if the ID is provided
	// if updatedTodo.ID == 0 {
	// 	return nil, fmt.Errorf("Invalid ID")
	// }
	s.Cache.Del(CacheKeyTodosList) // Invalidate cache
<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/service/service.go
>>>>>>> 4fe2e1c (unit tests)
========
	// Check if the ID is provided
	if updatedTodo.ID == 0 {
		return nil, fmt.Errorf("Invalid ID")
	}
	s.Cache.Del(constants.CacheKeyTodosList) // Invalidate cache
>>>>>>>> a07f97f (worker with refactoring code):pkg/service/service.go
<<<<<<< HEAD
=======
========
>>>>>>>> 4fe2e1c (unit tests):todo/service.go
>>>>>>> 4fe2e1c (unit tests)
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
<<<<<<< HEAD
	s.Cache.Del(constants.CacheKeyTodosList)
=======
	s.Cache.Del(CacheKeyTodosList)
>>>>>>> 4fe2e1c (unit tests)
	return nil
}

// Service function to cache todos
<<<<<<< HEAD
func (s *todoService) cacheTodos(todos []repo.Todo, cacheKey string) error {
=======
func (s *todoService) cacheTodos(todos []Todo, cacheKey string) error {
>>>>>>> 4fe2e1c (unit tests)
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
<<<<<<< HEAD
func (s *todoService) retrieveTodosFromCache(cacheKey string) ([]repo.Todo, error) {
=======
func (s *todoService) retrieveTodosFromCache(cacheKey string) ([]Todo, error) {
>>>>>>> 4fe2e1c (unit tests)
	// Attempt to retrieve todos from Redis cache
	cachedTodos, err := s.Cache.Get(cacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache from Redis: %w", err)
	}
	// if cachedTodos == "" {
	// 	return nil, nil // Cache miss
	// }

<<<<<<< HEAD
	var todos []repo.Todo
=======
	var todos []Todo
>>>>>>> 4fe2e1c (unit tests)
	// Deserialize JSON from cache into todos slice
	err = json.Unmarshal([]byte(cachedTodos), &todos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached todos: %w", err)
	}
	// fmt.Println("retreived from cache")
	return todos, nil
}
