package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/todo", createTodo).Methods("POST")
	r.HandleFunc("/todo", getTodos).Methods("GET")
	r.HandleFunc("/todo/{id}", updateTodo).Methods("PATCH")
	r.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
