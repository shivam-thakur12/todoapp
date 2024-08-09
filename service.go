package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

// Service function for creating a todo
func createTodoService(r *http.Request) (*Todo, error) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

// Service function for retrieving todos
func getTodosService() ([]Todo, error) {
	todos, err := getTodosRepo()
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// Service function for updating a todo
func updateTodoService(r *http.Request) (int, *Todo, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, nil, err
	}

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		return 0, nil, err
	}

	return id, &updatedTodo, nil
}

// Service function for deleting a todo
func deleteTodoService(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, err
	}
	return id, nil
}
