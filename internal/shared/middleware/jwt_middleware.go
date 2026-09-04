package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hurbbiee/todo-list-backend/internal/shared/auth"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
)

func JWTMiddleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(enum.AuthUnauthorized)
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Unauthorized(enum.AuthUnauthorized)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.ParseToken(tokenStr, secret)
		if err != nil {

			if errors.Is(err, jwt.ErrTokenExpired) {
				return response.Unauthorized(enum.AuthTokenExpired)
			}
			return response.Unauthorized(enum.AuthUnauthorized)
		}

		c.Locals("userId", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
