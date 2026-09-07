package notification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/config"
	delivery "github.com/hurbbiee/todo-list-backend/internal/modules/notification/delivery/http"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/repository"
	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(
	app *fiber.App,
	db *pgxpool.Pool,
	cfg *config.Config,
	cipher service.SecretCipher,
) {
	repo := repository.NewDiscordConnectionRepoPg(db)
	svc := service.NewDiscordConnectionService(repo, cipher)
	handler := delivery.NewDiscordConnectionHandler(svc)
	delivery.RegisterDiscordConnectionRoutes(app, handler, cfg.JWTSecret)
}
