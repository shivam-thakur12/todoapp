package routes

import (
	"TODO/pkg/handlers"

	"github.com/gorilla/mux"
)

// setupRoutes configures and returns a new mux.Router with defined routes.
func SetupRoutes(handler *handlers.TodoHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/todo", handler.Create).Methods("POST")
	r.HandleFunc("/todo", handler.Get).Methods("GET")
	r.HandleFunc("/todo/{id}", handler.Update).Methods("PATCH")
	r.HandleFunc("/todo/{id}", handler.Delete).Methods("DELETE")

	return r
}
