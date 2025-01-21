package todo

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /todos", handler.Create)
	mux.HandleFunc("GET /todos", handler.Get)
}
