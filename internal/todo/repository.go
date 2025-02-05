package todo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(todo *Todo) error {
	query := `INSERT INTO todos (title, description, completed, created_at, updated_at, user_id) VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, todo.Title, todo.Description, todo.Completed, now, now, todo.UserId)
	if err != nil {
		return fmt.Errorf("error creating todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	todo.Id = id
	todo.CreatedAt = now
	todo.UpdatedAt = now

	return nil
}

func (r *Repository) GetOneById(todoId int64, userId int64) (*Todo, error) {
	var todo Todo
	query := `SELECT id, title, description, completed, created_at, updated_at, user_id FROM todos WHERE id = ? AND user_id = ?`

	err := r.db.QueryRow(query, todoId, userId).Scan(&todo.Id,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("error fetching todo: %w", err)
	}

	return &todo, nil
}

func (r *Repository) UpdateCompleted(todoId int64, userId int64, completed bool) error {
	query := `UPDATE todos SET completed = ? WHERE id = ? AND user_id = ?`

	result, err := r.db.Exec(query, completed, todoId, userId)
	if err != nil {
		return fmt.Errorf("error updating todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("todo not found or user doesn't have permission to update")
	}

	return nil
}
