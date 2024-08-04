package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"todo/models"
)

var toDos = []models.ToDo{
	{Id: "1", Title: "First", Completed: false},
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getToDos(w, r)
	case "POST":
		createToDo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")
	switch r.Method {
	case "GET":
		getToDoById(w, r, id)
	case "PUT":
		updateToDo(w, r, id)
	case "DELETE":
		deleteToDo(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func getToDos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toDos)
}

func createToDo(w http.ResponseWriter, r *http.Request) {
	var newToDo models.ToDo
	if err := json.NewDecoder(r.Body).Decode(&newToDo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	toDos = append(toDos, newToDo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newToDo)
}

func getToDoById(w http.ResponseWriter, r *http.Request, id string) {
	for _, a := range toDos {
		if a.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(a)
			return
		}
	}

	http.Error(w, "To-Do not found", http.StatusNotFound)
}

func updateToDo(w http.ResponseWriter, r *http.Request, id string) {
	var updatedToDo models.ToDo
	if err := json.NewDecoder(r.Body).Decode(&updatedToDo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, t := range toDos {
		if t.Id == id {
			toDos[i] = updatedToDo
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedToDo)
			return
		}
	}
	http.Error(w, "To-do not found", http.StatusNotFound)
}

func deleteToDo(w http.ResponseWriter, r *http.Request, id string) {
	for i, t := range toDos {
		if t.Id == id {
			toDos = append(toDos[:i], toDos[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "To-do deleted"})
			return
		}
	}
	http.Error(w, "To-do not found", http.StatusNotFound)
}
