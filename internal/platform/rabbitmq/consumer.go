package rabbitmq

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TodoCreatedQueue      = "notification.todo.created"
	TodoCreatedRoutingKey = "todo.created"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(
	connection *Connection,
) (*Consumer, error) {
	if connection == nil || connection.IsClosed() {
		return nil, errors.New(
			"rabbitmq connection is not available",
		)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	// Consumer ต้อง declare exchange เองด้วย
	// เพื่อให้ Worker เปิดก่อน API ได้
	err = channel.ExchangeDeclare(
		TodoExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf(
			"declare exchange %s: %w",
			TodoExchange,
			err,
		)
	}

	queue, err := channel.QueueDeclare(
		TodoCreatedQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf(
			"declare queue %s: %w",
			TodoCreatedQueue,
			err,
		)
	}

	err = channel.QueueBind(
		queue.Name,
		TodoCreatedRoutingKey,
		TodoExchange,
		false,
		nil,
	)
	if err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf(
			"bind queue %s to exchange %s: %w",
			queue.Name,
			TodoExchange,
			err,
		)
	}

	// Worker รับทีละหนึ่ง message
	// และต้อง ACK message เดิมก่อนรับงานถัดไป
	err = channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf(
			"set rabbitmq qos: %w",
			err,
		)
	}

	return &Consumer{
		channel: channel,
	}, nil
}

func (c *Consumer) ConsumeTodoCreated() (
	<-chan amqp.Delivery,
	error,
) {
	if c == nil || c.channel == nil {
		return nil, errors.New(
			"rabbitmq consumer is not initialized",
		)
	}

	deliveries, err := c.channel.Consume(
		TodoCreatedQueue,
		"todo-created-notification-worker",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"consume queue %s: %w",
			TodoCreatedQueue,
			err,
		)
	}

	return deliveries, nil
}

func (c *Consumer) Close() error {
	if c == nil || c.channel == nil {
		return nil
	}

	if c.channel.IsClosed() {
		return nil
	}

	if err := c.channel.Close(); err != nil {
		return fmt.Errorf(
			"close rabbitmq consumer: %w",
			err,
		)
	}

	return nil
}
