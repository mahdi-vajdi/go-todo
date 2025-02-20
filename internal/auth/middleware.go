package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const UserContextKey string = "user"

func JwtMiddleware(authService *Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			// Check if auth header is valid
			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			// Remove "Bearer " prefix
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			userId, err := authService.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			// Store user information in context
			c.Set(UserContextKey, userId)

			return next(c)
		}
	}
}
