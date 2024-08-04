package main

import (
	"log"
	"net/http"
	"todo/handlers"
)

func main() {
	http.HandleFunc("/todos", handlers.TodosHandler)
	http.HandleFunc("/todos/", handlers.TodoHandler)

	log.Println("Server is running on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
