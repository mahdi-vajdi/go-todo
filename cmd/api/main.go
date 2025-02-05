package main

import (
	"log"
	"net/http"
	"todo/config"
	"todo/internal/auth"
	"todo/internal/database"
	"todo/internal/todo"
)

func main() {
	// Load the configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load cofig: %v", err)
	}

	// Initialize the database
	db, err := database.NewConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the repositories
	authRepository := auth.NewRepository(db)
	todoRepository := todo.NewRepository(db)

	// Initialize services
	authService := auth.NewService(authRepository, &cfg.Auth)
	todoService := todo.NewService(todoRepository)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)

	mux := http.NewServeMux()

	// Setup routes
	// Auth
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	// To-do
	mux.Handle("POST /todos", auth.JwtMiddleware(&cfg.Auth, http.HandlerFunc(todoHandler.Create)))
	mux.Handle("GET /todos", auth.JwtMiddleware(&cfg.Auth, http.HandlerFunc(todoHandler.Get)))
	mux.Handle("POST /todos/done", auth.JwtMiddleware(&cfg.Auth, http.HandlerFunc(todoHandler.SetDone)))

	// Add global prefix to the routes
	apiHandler := http.StripPrefix("/api", mux)

	// Start the server
	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", apiHandler); err != nil {
		log.Fatal(err)
	}
}
