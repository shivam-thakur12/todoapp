<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/repo/repo.go
>>>>>>> 4fe2e1c (unit tests)
<<<<<<<< HEAD:todo/repo.go
package todo
========
package repo

<<<<<<< HEAD
import "TODO/pkg/server"
>>>>>>>> a07f97f (worker with refactoring code):pkg/repo/repo.go
=======
<<<<<<< HEAD:pkg/repo/repo.go
import "TODO/pkg/server"
>>>>>>>> a07f97f (worker with refactoring code):pkg/repo/repo.go
=======
import "TODO/todo/server"
>>>>>>> 8bb23e4 (more refactoring done.):todo/repo/repo.go
========
package todo
>>>>>>>> 4fe2e1c (unit tests):todo/repo.go
>>>>>>> 4fe2e1c (unit tests)

type Todo struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type TodoRepository interface {
	CreateTodoRepo(todo *Todo) error
	GetTodosRepo() ([]Todo, error)
	UpdateTodoRepo(todo *Todo) (int64, error)
	// DeleteTodoRepo(id int) (int64, error)
}

type todoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) CreateTodoRepo(todo *Todo) error {
<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/repo/repo.go
>>>>>>> 4fe2e1c (unit tests)
<<<<<<<< HEAD:todo/repo.go
	err := DB.QueryRow("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id", todo.Title, todo.Status).Scan(&todo.ID)
========
	err := server.DB.QueryRow("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id", todo.Title, todo.Status).Scan(&todo.ID)
>>>>>>>> a07f97f (worker with refactoring code):pkg/repo/repo.go
<<<<<<< HEAD
=======
========
	err := DB.QueryRow("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id", todo.Title, todo.Status).Scan(&todo.ID)
>>>>>>>> 4fe2e1c (unit tests):todo/repo.go
>>>>>>> 4fe2e1c (unit tests)
	if err != nil {
		return err
	}
	return nil
}

// Repository function for retrieving todos
func (r *todoRepository) GetTodosRepo() ([]Todo, error) {
<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/repo/repo.go
>>>>>>> 4fe2e1c (unit tests)
<<<<<<<< HEAD:todo/repo.go
	rows, err := DB.Query("SELECT id, title, status FROM todos WHERE deleted_at IS NULL")
========
	rows, err := server.DB.Query("SELECT id, title, status FROM todos WHERE deleted_at IS NULL")
>>>>>>>> a07f97f (worker with refactoring code):pkg/repo/repo.go
<<<<<<< HEAD
=======
========
	rows, err := DB.Query("SELECT id, title, status FROM todos WHERE deleted_at IS NULL")
>>>>>>>> 4fe2e1c (unit tests):todo/repo.go
>>>>>>> 4fe2e1c (unit tests)
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
<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/repo/repo.go
>>>>>>> 4fe2e1c (unit tests)
<<<<<<<< HEAD:todo/repo.go
	res, err := DB.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3 AND deleted_at IS NULL", todo.Title, todo.Status, todo.ID)
========
	res, err := server.DB.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3 AND deleted_at IS NULL", todo.Title, todo.Status, todo.ID)
>>>>>>>> a07f97f (worker with refactoring code):pkg/repo/repo.go
<<<<<<< HEAD
=======
========
	res, err := DB.Exec("UPDATE todos SET title=$1, status=$2 WHERE id=$3 AND deleted_at IS NULL", todo.Title, todo.Status, todo.ID)
>>>>>>>> 4fe2e1c (unit tests):todo/repo.go
>>>>>>> 4fe2e1c (unit tests)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

<<<<<<< HEAD
<<<<<<<< HEAD:todo/repo.go
=======
<<<<<<<< HEAD:todo/repo/repo.go
<<<<<<<< HEAD:todo/repo.go
========
>>>>>>>> 4fe2e1c (unit tests):todo/repo.go
>>>>>>> 4fe2e1c (unit tests)
// // Repository function for deleting a todo
// func (r *todoRepository) DeleteTodoRepo(id int) (int64, error) {
// 	res, err := DB.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
// 	if err != nil {
// 		return 0, err
// 	}
<<<<<<< HEAD
=======
<<<<<<<< HEAD:todo/repo/repo.go
>>>>>>> 4fe2e1c (unit tests)
========
// Repository function for deleting a todo
func (r *todoRepository) DeleteTodoRepo(id int) (int64, error) {
	res, err := server.DB.Exec("UPDATE todos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, err
	}
>>>>>>>> a07f97f (worker with refactoring code):pkg/repo/repo.go
<<<<<<< HEAD
=======
========
>>>>>>>> 4fe2e1c (unit tests):todo/repo.go
>>>>>>> 4fe2e1c (unit tests)

// 	rowsAffected, err := res.RowsAffected()
// 	return rowsAffected, err
// }
