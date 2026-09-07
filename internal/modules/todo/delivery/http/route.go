package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/shared/middleware"
)

func RegisterTodoRoutes(
	app *fiber.App,
	handler *TodoHandler,
	jwtSecret []byte,
) {
	api := app.Group("/api", middleware.JWTMiddleware(jwtSecret))

	todo := api.Group("/todo")
	todo.Post("/search", handler.Search)
	todo.Post("/create", handler.Create)
	todo.Get("/counts", handler.CountStatus)
	todo.Patch("/:id/update", handler.Update)
	todo.Patch("/:id/update-complete", handler.UpdateStatus)
	todo.Delete("/:id/delete", handler.Delete)
}
