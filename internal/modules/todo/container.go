package todo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/config"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/delivery/http"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/repository"
	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(
	app *fiber.App,
	db *pgxpool.Pool,
	cfg *config.Config,
	publisher service.EventPublisher,
) {
	repo := repository.NewTodoRepopg(db)
	svc := service.NewTodoService(repo, publisher)
	handler := http.NewTodoHandler(svc)
	http.RegisterTodoRoutes(app, handler, cfg.JWTSecret)
}
