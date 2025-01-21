package main

import (
	"log"
	"net/http"
	"todo/config"
	"todo/internal/auth"
	"todo/internal/database"
	router "todo/internal/http"
	"todo/internal/todo"
)

func main() {
	// Load the configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load cofig: %v", err)
	}

	// Initialize the database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
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
	mux := router.NewRouter(authHandler, todoHandler)

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
