package main

import (
	"github.com/gorilla/mux"
)

// setupRoutes configures and returns a new mux.Router with defined routes.
func setupRoutes(handler *TodoHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/todo", handler.createTodo).Methods("POST")
	r.HandleFunc("/todo", handler.getTodos).Methods("GET")
	r.HandleFunc("/todo/{id}", handler.updateTodo).Methods("PATCH")
	r.HandleFunc("/todo/{id}", handler.deleteTodo).Methods("DELETE")

	return r
}
