package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"todo/models"
	"todo/repository"
)

var repo *repository.TodoRepository

func InitHandlers(db *sql.DB) {
	repo = &repository.TodoRepository{DB: db}
}

func GetToDos(w http.ResponseWriter, r *http.Request) {
	todos, err := repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateToDo(w http.ResponseWriter, r *http.Request) {
	var newToDo models.ToDo
	if err := json.NewDecoder(r.Body).Decode(&newToDo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repo.Create(&newToDo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newToDo)
}

func GetToDoById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := repo.GetOneById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todo == nil {
		http.Error(w, "To-Do not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid Id",
			http.StatusBadRequest)
		return
	}

	var updatedToDo models.ToDo
	if err := json.NewDecoder(r.Body).Decode(&updatedToDo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedToDo.Id = id

	if err := repo.Update(&updatedToDo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedToDo)
}

func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid to-do ID", http.StatusBadRequest)
		return
	}

	if err := repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "To-Do deleted"})
}
