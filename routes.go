package main

import (
	"github.com/gorilla/mux"
)

// setupRoutes configures and returns a new mux.Router with defined routes.
func setupRoutes(handler *TodoHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/todo", handler.create).Methods("POST")
	r.HandleFunc("/todo", handler.get).Methods("GET")
	r.HandleFunc("/todo/{id}", handler.update).Methods("PATCH")
	r.HandleFunc("/todo/{id}", handler.delete).Methods("DELETE")

	return r
}
