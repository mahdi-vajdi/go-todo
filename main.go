package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo/database"
	"todo/handlers"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()
	handlers.InitHandlers(database.DB)

	r.HandleFunc("/todos", handlers.GetToDos).Methods("GET")
	r.HandleFunc("/todos", handlers.CreateToDo).Methods("POST")
	r.HandleFunc("/todos/{id}", handlers.GetToDoById).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.UpdateToDo).Methods("PUT")
	r.HandleFunc("/todos/{id}", handlers.DeleteToDo).Methods("DELETE")

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
