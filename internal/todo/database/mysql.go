package todo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"todo/internal/todo"
)

type MysqlRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) todo.Repository {
	return &MysqlRepository{db: db}
}

func (r *MysqlRepository) Create(todo *todo.Todo) error {
	query := `INSERT INTO todos (title, description, completed, created_at, updated_at, user_id) VALUES (:title, :description, :completed, :created_at, :updated_at, :user_id)`

	_, err := r.db.NamedExec(query, map[string]interface{}{
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
		"created_at":  todo.CreatedAt,
		"updated_at":  todo.UpdatedAt,
		"user_id":     todo.UserId,
	})
	if err != nil {
		return fmt.Errorf("error creating todo: %w", err)
	}

	return nil
}

func (r *MysqlRepository) GetOneById(todoId int64, userId int64) (*todo.Todo, error) {
	var t todo.Todo
	query := `SELECT id, title, description, completed, created_at, updated_at, user_id 
			  FROM todos 
			  WHERE id = ? AND user_id = ?`

	err := r.db.Get(&t, query, todoId, userId)
	if err != nil {
		if err == sql.ErrNoRows { // Don't use errors.Is
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("error fetching todo: %w", err)
	}

	return &t, nil
}

func (r *MysqlRepository) UpdateCompleted(todoId int64, userId int64, completed bool) error {
	query := `UPDATE todos SET completed = :completed WHERE id = :todo_id AND user_id = :user_id`

	_, err := r.db.NamedExec(query, map[string]interface{}{
		"completed": completed,
		"todo_id":   todoId,
		"user_id":   userId,
	})
	if err != nil {
		return fmt.Errorf("error updating todo: %w", err)
	}

	return nil
}
