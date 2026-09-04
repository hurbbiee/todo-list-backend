package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/config"
	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/delivery/http"
	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/repository"
	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(
	app *fiber.App,
	db *pgxpool.Pool,
	cfg *config.Config,
) {
	repo := repository.NewAuthRepoPg(db)
	svc := service.NewAuthService(repo, cfg.JWTSecret)
	handler := http.NewAuthHandler(svc)
	http.RegisterAuthRoutes(app, handler, cfg.JWTSecret)
}
