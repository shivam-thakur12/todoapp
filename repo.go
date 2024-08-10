package main

type TodoRepository interface {
	CreateTodoRepo(todo *Todo) error
	GetTodosRepo() ([]Todo, error)
	UpdateTodoRepo(todo *Todo) (int64, error)
	DeleteTodoRepo(id int) (int64, error)
}

type todoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) CreateTodoRepo(todo *Todo) error {
	err := db.QueryRow("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id", todo.Title, todo.Status).Scan(&todo.ID)
	if err != nil {
		return err
	}
	return nil
}

// Repository function for retrieving todos
func (r *todoRepository) GetTodosRepo() ([]Todo, error) {
	rows, err := db.Query("SELECT id, title, status FROM todos WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// Repository function for updating a todo
func (r *todoRepository) UpdateTodoRepo(todo *Todo) (int64, error) {
	res, err := db.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3 AND deleted_at IS NULL", todo.Title, todo.Status, todo.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Repository function for deleting a todo
func (r *todoRepository) DeleteTodoRepo(id int) (int64, error) {
	res, err := db.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err
}
