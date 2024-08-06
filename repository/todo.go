package repository

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"todo/models"
)

type TodoRepository struct {
	DB *sql.DB
}

func (r *TodoRepository) GetAll() ([]models.ToDo, error) {
	rows, err := r.DB.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.ToDo
	for rows.Next() {
		var todo models.ToDo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *TodoRepository) Create(todo *models.ToDo) error {
	result, err := r.DB.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", todo.Title, todo.Completed)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	todo.Id = int(id)
	return nil
}

func (r *TodoRepository) GetOneById(id int) (*models.ToDo, error) {
	var todo models.ToDo

	err := r.DB.QueryRow("SELECT id, title, completed FROM todos WHER id = ?", id).Scan(&todo.Id, &todo.Title, &todo.Completed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(todo *models.ToDo) error {
	_, err := r.DB.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ?", todo.Title, todo.Completed, todo.Id)
	return err
}

func (r *TodoRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
