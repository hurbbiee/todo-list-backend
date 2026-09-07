package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TodoExchange = "todo.events"
)

type Publisher struct {
	channel *amqp.Channel
	mu      sync.Mutex
}

func NewPublisher(
	connection *Connection,
) (*Publisher, error) {
	if connection == nil || connection.IsClosed() {
		return nil, errors.New(
			"rabbitmq connection is not available",
		)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

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

	return &Publisher{
		channel: channel,
	}, nil
}

func (p *Publisher) PublishJSON(
	ctx context.Context,
	routingKey string,
	payload any,
) error {
	if p == nil || p.channel == nil {
		return errors.New(
			"rabbitmq publisher is not initialized",
		)
	}

	if routingKey == "" {
		return errors.New(
			"rabbitmq routing key is required",
		)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf(
			"marshal rabbitmq payload: %w",
			err,
		)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	err = p.channel.PublishWithContext(
		ctx,
		TodoExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Type:         routingKey,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"publish rabbitmq message with routing key %s: %w",
			routingKey,
			err,
		)
	}

	return nil
}

func (p *Publisher) Close() error {
	if p == nil || p.channel == nil {
		return nil
	}

	if p.channel.IsClosed() {
		return nil
	}

	if err := p.channel.Close(); err != nil {
		return fmt.Errorf(
			"close rabbitmq publisher: %w",
			err,
		)
	}

	return nil
}
