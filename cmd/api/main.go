package main

import (
	"log"
	"net/http"
	"todo/internal/auth"
	"todo/internal/platform/config"
	"todo/internal/platform/database"
	"todo/internal/todo"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load cofig: %v", err)
	}

	// Initialize the database
	dbConfig := database.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Name:     cfg.Name,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize services
	authService := auth.NewService(db, []byte("secret"))
	todoService := todo.NewService(db)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// To-do Routes
	todoRoutes := http.NewServeMux()
	todoRoutes.HandleFunc("/", todoHandler.Create)
	todoRoutes.HandleFunc("/get", todoHandler.Get)

	// Apply JWT middleware
	mux.Handle("/todos", auth.JwtMiddleware(http.StripPrefix("/todos", todoRoutes)))
	mux.Handle("/todos/", auth.JwtMiddleware(http.StripPrefix("/todos", todoRoutes)))

	// Add global prefix to the routes
	apiHandler := http.StripPrefix("/api", mux)

	// Start the server
	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", apiHandler); err != nil {
		log.Fatal(err)
	}
}
