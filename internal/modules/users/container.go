package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/config"
	"github.com/hurbbiee/todo-list-backend/internal/modules/users/delivery/http"
	"github.com/hurbbiee/todo-list-backend/internal/modules/users/repository"
	"github.com/hurbbiee/todo-list-backend/internal/modules/users/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(
	app *fiber.App,
	db *pgxpool.Pool,
	cfg *config.Config,
) {
	repo := repository.NewUserRepoPg(db)
	svc := service.NewUserService(repo)
	handler := http.NewUserHandler(svc)
	http.RegisterUserRoutes(app, handler, cfg.JWTSecret)
}
