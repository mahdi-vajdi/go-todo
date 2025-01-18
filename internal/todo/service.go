package todo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(todo *Todo) error {
	query := `INSERT INTO todos (title, description, completed, created_at, updated_at, user_id) VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := s.db.Exec(query, todo.Title, todo.Description, todo.Completed, now, now, todo.UserId)
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

func (s *Service) GetById(todoId int64, userId int64) (*Todo, error) {
	var todo Todo
	query := `SELECT id, title, description, completed, created_at, updated_at, user_id FROM todos WHERE id = ? AND user_id = ?`

	err := s.db.QueryRow(query, todoId, userId).Scan(&todo.Id,
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
