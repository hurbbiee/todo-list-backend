package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hurbbiee/todo-list-backend/internal/config"
	authRegister "github.com/hurbbiee/todo-list-backend/internal/modules/auth"
	todoRegister "github.com/hurbbiee/todo-list-backend/internal/modules/todo"
	userRegister "github.com/hurbbiee/todo-list-backend/internal/modules/user"
	db "github.com/hurbbiee/todo-list-backend/internal/platform/database"
	rabbitmqClient "github.com/hurbbiee/todo-list-backend/internal/platform/rabbitmq"
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

	rabbitConnection, err := rabbitmqClient.NewConnection(
		cfg.RabbitMQ.URL,
	)
	if err != nil {
		log.Fatal("connect rabbitmq: ", err)
	}
	defer func() {
		if err := rabbitConnection.Close(); err != nil {
			log.Printf("close rabbitmq: %v", err)
		}
	}()

	log.Println("rabbitmq connected")

	rabbitPublisher, err := rabbitmqClient.NewPublisher(
		rabbitConnection,
	)
	if err != nil {
		log.Printf(
			"create rabbitmq publisher: %v",
			err,
		)
		return
	}

	defer func() {
		if err := rabbitPublisher.Close(); err != nil {
			log.Printf(
				"close rabbitmq publisher: %v",
				err,
			)
		}
	}()

	log.Println("rabbitmq publisher ready")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			status, resp := response.FromError(err)
			return c.Status(status).JSON(resp)
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.FrontendURL,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
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
	todoRegister.Register(app, databasePool, &cfg, rabbitPublisher)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Printf("fiber server stopped: %v", err)
	}
}
