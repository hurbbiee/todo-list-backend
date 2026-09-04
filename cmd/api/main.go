package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hurbbiee/todo-list-backend/internal/config"
	authRegister "github.com/hurbbiee/todo-list-backend/internal/modules/auth"
	userRegister "github.com/hurbbiee/todo-list-backend/internal/modules/users"
	db "github.com/hurbbiee/todo-list-backend/internal/platform/database"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
)

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("Failed to load timezone Asia/Bangkok: %v", err)
		loc = time.FixedZone("BKK", 7*60*60)
	}

	time.Local = loc
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load config: ", err)
	}

	databasePool, err := db.NewPostgres(
		context.Background(),
		cfg.Database,
	)
	if err != nil {
		log.Fatal("connect database: ", err)
	}
	defer databasePool.Close()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			status, resp := response.FromError(err)
			return c.Status(status).JSON(resp)
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.FrontendURL,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	authRegister.Register(app, databasePool, &cfg)
	userRegister.Register(app, databasePool, &cfg)

	app.Listen(":" + cfg.Port)
}
