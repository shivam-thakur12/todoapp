package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupRoutes configures and returns a new mux.Router with defined routes.
func SetupRoutes(handler *TodoHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/todo", handler.Create).Methods("POST")
	r.HandleFunc("/todo", handler.Get).Methods("GET")
	r.HandleFunc("/todo/{id}", handler.Update).Methods("PATCH")
	r.HandleFunc("/todo/{id}", handler.Delete).Methods("DELETE")
	r.Path("/metrics").Handler(promhttp.Handler())
	return r
}
