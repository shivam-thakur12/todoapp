package main

import (
	"encoding/json"
	"net/http"
)

// Handle creating a new todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	todo, err := createTodoService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = createTodoRepo(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// Handle retrieving all todos
func getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := getTodosService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Handle updating a todo by ID
func updateTodo(w http.ResponseWriter, r *http.Request) {
	id, updatedTodo, err := updateTodoService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rowsAffected, err := updateTodoRepo(id, updatedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	updatedTodo.ID = id
	json.NewEncoder(w).Encode(updatedTodo)
}

// Handle deleting a todo by ID
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := deleteTodoService(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := deleteTodoRepo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
