package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Config struct {
	Database DatabaseConfig `toml:"database"`
}

type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
	Sslmode  string `toml:"sslmode"`
}

var db *sql.DB

type Todo struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	DeletedAt *string `json:"deleted_at,omitempty"` // Pointer to handle NULL values
}

// Initialize the database connection
func initDB() {
	var config Config

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Dbname, config.Database.Sslmode)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected successfully!")
}

// Handle creating a new todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	//json.NewDecoder(r.Body).Decode(&todo) creates a new JSON decoder and reads from request body
	// and tries to decodes the JSON data into the todo var. if error, displayed
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id", todo.Title, todo.Status).Scan(&todo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// Handle retrieving all todos
func getTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, status FROM todos WHERE deleted_at IS NULL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Handle updating a todo by ID
func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3 AND deleted_at IS NULL", updatedTodo.Title, updatedTodo.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Set deleted_at to the current time
	res, err := db.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
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

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/todo", createTodo).Methods("POST")
	r.HandleFunc("/todo", getTodos).Methods("GET")
	r.HandleFunc("/todo/{id}", updateTodo).Methods("PATCH")
	r.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
