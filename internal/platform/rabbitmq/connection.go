package rabbitmq

import (
	"errors"
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	connection *amqp.Connection
}

func NewConnection(url string) (*Connection, error) {
	if strings.TrimSpace(url) == "" {
		return nil, errors.New(
			"rabbitmq url is required",
		)
	}

	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf(
			"connect rabbitmq: %w",
			err,
		)
	}

	return &Connection{
		connection: connection,
	}, nil
}

func (c *Connection) Channel() (*amqp.Channel, error) {
	if c == nil || c.connection == nil {
		return nil, errors.New(
			"rabbitmq connection is not initialized",
		)
	}

	if c.connection.IsClosed() {
		return nil, errors.New(
			"rabbitmq connection is closed",
		)
	}

	channel, err := c.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf(
			"open rabbitmq channel: %w",
			err,
		)
	}

	return channel, nil
}

func (c *Connection) IsClosed() bool {
	if c == nil || c.connection == nil {
		return true
	}

	return c.connection.IsClosed()
}

func (c *Connection) Close() error {
	if c == nil || c.connection == nil {
		return nil
	}

	if c.connection.IsClosed() {
		return nil
	}

	if err := c.connection.Close(); err != nil {
		return fmt.Errorf(
			"close rabbitmq connection: %w",
			err,
		)
	}

	return nil
}
