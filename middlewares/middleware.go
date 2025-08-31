package middlewares

import (
	"net/http"
	"strings"

	"billing-engine/global"
	"billing-engine/utils/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type GoMiddlewareInterface interface {
	ValidateCORS(next echo.HandlerFunc) echo.HandlerFunc
	ValidateToken(next echo.HandlerFunc) echo.HandlerFunc
}

type GoMiddleware struct {
	AccessSecret []byte
}

func InitMiddleware(accessSecret []byte) GoMiddlewareInterface {
	return &GoMiddleware{AccessSecret: accessSecret}
}

func (gm *GoMiddleware) ValidateCORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

func (gm *GoMiddleware) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, global.BadResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization header is missing",
			})
		}

		parts := strings.Split(authHeader, " ")
		if (len(parts) != 2) || strings.ToLower(parts[0]) != "bearer" {
			return c.JSON(http.StatusUnauthorized, global.BadResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Authorization header format",
			})
		}
		accessToken := parts[1]
		claims := &token.Claims{}
		token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
			return gm.AccessSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, global.BadResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid access token",
			})
		}
		return next(c)
	}
}
