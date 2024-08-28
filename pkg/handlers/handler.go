<<<<<<<< HEAD:cmd/api/handler.go
package api

import (
	"TODO/todo"
========
package handlers

import (
	"TODO/pkg/service"
>>>>>>>> a07f97f (worker with refactoring code):pkg/handlers/handler.go
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
<<<<<<<< HEAD:cmd/api/handler.go
	Service todo.TodoService
========
	Service service.TodoService
>>>>>>>> a07f97f (worker with refactoring code):pkg/handlers/handler.go
}

// Handle creating a new todo
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	todo, err := h.Service.CreateTodoService(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// Handle retrieving all todos
func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	todos, err := h.Service.GetTodosService()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	updatedTodo, err := h.Service.UpdateTodoService(id, r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(updatedTodo)
}

// Handle deleting a todo by ID
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.Service.DeleteTodoService(id, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
