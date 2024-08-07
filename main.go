package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Config struct {
	Database   DatabaseConfig  `toml:"database"`
	Migrations MigrationConfig `toml:"migrations"`
}

type MigrationConfig struct {
	Path string `toml:"path"`
}

type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
	Sslmode  string `toml:"sslmode"`
}

var db *sql.DB

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// Initialize the database connection and run migrations
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

	config.runMigrations()

}
func (config Config) runMigrations() {
	// Adjust connection string format for migrations
	migrationConnStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Dbname, config.Database.Sslmode)

	// Run migrations
	// Connection string to the database
	m, err := migrate.New(
		fmt.Sprintf("file://%s", config.Migrations.Path), // Source path to migration files from config
		migrationConnStr)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("Database migrated successfully!")
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
	rows, err := db.Query("SELECT id, title, status FROM todos")
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

	res, err := db.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3", updatedTodo.Title, updatedTodo.Status, id)
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

	res, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
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
