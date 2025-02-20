package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("error closing database connectio %s", err.Error())
		}
	}(db)

	// Initialize the repositories
	authRepository := auth.NewRepository(db)
	todoRepository := todo.NewRepository(db)

	// Initialize the services
	authService := auth.NewService(authRepository, &cfg.Auth)
	todoService := todo.NewService(todoRepository)

	// Initialize the handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)

	// Create the http server
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup apiRoute prefix
	apiRoute := e.Group("/api")

	// Auth routes
	apiRoute.POST("/register", authHandler.Register)
	apiRoute.POST("/login", authHandler.Login)

	// To-do routes
	todoRoute := apiRoute.Group("/todos")
	todoRoute.Use(auth.JwtMiddleware(authService))
	todoRoute.POST("", todoHandler.Create)
	todoRoute.GET("/:id", todoHandler.Get)
	todoRoute.PATCH("/:id", todoHandler.Update)

	// Handle not found routes
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found"})
	})

	// Start the server
	log.Print("server starting on port 8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
