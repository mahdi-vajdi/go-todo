package todo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"todo/internal/todo"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) todo.Repository {
	return &MysqlRepository{db: db}
}

func (r *MysqlRepository) Create(todo *todo.Todo) error {
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

func (r *MysqlRepository) GetOneById(todoId int64, userId int64) (*todo.Todo, error) {
	var t todo.Todo
	query := `SELECT id, title, description, completed, created_at, updated_at, user_id FROM todos WHERE id = ? AND user_id = ?`

	err := r.db.QueryRow(query, todoId, userId).Scan(&t.Id,
		&t.Title,
		&t.Description,
		&t.Completed,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("t not found")
		}
		return nil, fmt.Errorf("error fetching t: %w", err)
	}

	return &t, nil
}

func (r *MysqlRepository) UpdateCompleted(todoId int64, userId int64, completed bool) error {
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
