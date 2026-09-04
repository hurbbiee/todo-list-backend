package http

import "github.com/gofiber/fiber/v2"

func RegisterAuthRoutes(
	app *fiber.App,
	handler *AuthHandler,
	jwtSecret []byte,
) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
}