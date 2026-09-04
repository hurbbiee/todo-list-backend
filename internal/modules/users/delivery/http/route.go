package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/shared/middleware"
)

func RegisterUserRoutes(
	app *fiber.App,
	handler *UserHandler,
	jwtSecret []byte,
) {
	api := app.Group("/api", middleware.JWTMiddleware(jwtSecret))

	user := api.Group("/user")
	user.Post("/create", handler.Create)
	user.Get("/profile", handler.Profile)
}
