package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/shared/middleware"
)

func RegisterDiscordConnectionRoutes(
	app *fiber.App,
	handler *DiscordConnectionHandler,
	jwtSecret []byte,
) {
	api := app.Group("/api", middleware.JWTMiddleware(jwtSecret))
	discord := api.Group("/notification/discord")

	discord.Get("/", handler.Status)
	discord.Put("/", handler.Save)
	discord.Patch("/enabled", handler.UpdateEnabled)
	discord.Delete("/", handler.Disconnect)
}
