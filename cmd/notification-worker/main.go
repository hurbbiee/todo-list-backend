package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hurbbiee/todo-list-backend/internal/config"

	notificationConsumer "github.com/hurbbiee/todo-list-backend/internal/modules/notification/consumer"
	rabbitmqClient "github.com/hurbbiee/todo-list-backend/internal/platform/rabbitmq"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Printf(
			"load config: %v",
			err,
		)
		return
	}

	rabbitConnection, err := rabbitmqClient.NewConnection(
		cfg.RabbitMQ.URL,
	)
	if err != nil {
		log.Printf(
			"connect rabbitmq: %v",
			err,
		)
		return
	}

	defer func() {
		if err := rabbitConnection.Close(); err != nil {
			log.Printf(
				"close rabbitmq connection: %v",
				err,
			)
		}
	}()

	consumer, err := rabbitmqClient.NewConsumer(
		rabbitConnection,
	)
	if err != nil {
		log.Printf(
			"create rabbitmq consumer: %v",
			err,
		)
		return
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf(
				"close rabbitmq consumer: %v",
				err,
			)
		}
	}()

	deliveries, err := consumer.ConsumeTodoCreated()
	if err != nil {
		log.Printf(
			"consume todo.created: %v",
			err,
		)
		return
	}

	log.Println(
		"notification worker is waiting for todo.created messages",
	)

	todoCreatedHandler :=
		notificationConsumer.NewTodoCreatedHandler()
		
	log.Println(
		"notification worker is waiting for todo.created messages",
	)

	for {
		select {
		case <-ctx.Done():
			log.Println(
				"notification worker shutting down",
			)
			return

		case delivery, ok := <-deliveries:
			if !ok {
				log.Println(
					"rabbitmq delivery channel closed",
				)
				return
			}

			err := todoCreatedHandler.Handle(
				ctx,
				delivery.Body,
			)
			if err != nil {
				log.Printf(
					"handle todo.created failed: %v",
					err,
				)

				if errors.Is(
					err,
					notificationConsumer.ErrInvalidMessage,
				) {
					// Payload ผิด ต่อให้ส่งใหม่ก็ยังผิด
					// ไม่ requeue เพื่อป้องกัน infinite loop
					if rejectErr := delivery.Reject(false); rejectErr != nil {
						log.Printf(
							"reject rabbitmq message: %v",
							rejectErr,
						)
					}

					continue
				}

				// Error ชั่วคราว เช่น Discord ล่ม
				// ตอนนี้ส่งกลับ Queue ก่อน
				if nackErr := delivery.Nack(
					false,
					true,
				); nackErr != nil {
					log.Printf(
						"nack rabbitmq message: %v",
						nackErr,
					)
				}

				continue
			}

			if err := delivery.Ack(false); err != nil {
				log.Printf(
					"ack rabbitmq message: %v",
					err,
				)
			}
		}
	}
}
