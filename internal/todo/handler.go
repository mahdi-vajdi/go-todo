package todo

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"todo/internal/auth"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
	userId := c.Get(auth.UserContextKey).(int64)

	var todo Todo
	if err := c.Bind(&todo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	todo.UserId = userId

	if err := h.service.Create(&todo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Get(c echo.Context) error {
	userId := c.Get(auth.UserContextKey).(int64)
	todoId := c.Param("id")

	id, err := strconv.ParseInt(todoId, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID parameter")
	}

	todo, err := h.service.GetById(id, userId)
	if err != nil {
		if err.Error() == "todo not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *Handler) Update(c echo.Context) error {
	type requestBody struct {
		Completed bool `json:"completed"`
	}

	userId := c.Get(auth.UserContextKey).(int64)
	todoIdString := c.Param("id")
	var body requestBody

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	var todoId, err = strconv.ParseInt(todoIdString, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Todo ID")
	}

	err = h.service.UpdateCompleted(todoId, userId, body.Completed)
	if err != nil {
		if err.Error() == "todo not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Invalid Todo ID")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.NoContent(http.StatusOK)
}
