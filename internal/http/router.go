package http

import (
	"net/http"
	"todo/internal/auth"
	"todo/internal/todo"
)

func NewRouter(authHandler *auth.Handler, todoHandler *todo.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	auth.RegisterRoutes(mux, authHandler)
	todo.RegisterRoutes(mux, todoHandler)

	return mux
}
