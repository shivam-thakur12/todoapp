package main

// Repository function for creating a todo
func createTodoRepo(todo *Todo) error {
	err := db.QueryRow("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id", todo.Title, todo.Status).Scan(&todo.ID)
	return err
}

// Repository function for retrieving todos
func getTodosRepo() ([]Todo, error) {
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
func updateTodoRepo(id int, updatedTodo *Todo) (int64, error) {
	res, err := db.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3 AND deleted_at IS NULL", updatedTodo.Title, updatedTodo.Status, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err
}

// Repository function for deleting a todo
func deleteTodoRepo(id int) (int64, error) {
	res, err := db.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err
}
