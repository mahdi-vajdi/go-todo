package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c echo.Context) error {
	var credentials Credentials

	if err := c.Bind(&credentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.service.RegisterUser(credentials); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Login(c echo.Context) error {
	var credentials Credentials

	if err := c.Bind(&credentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.service.Login(credentials)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, response)
}
