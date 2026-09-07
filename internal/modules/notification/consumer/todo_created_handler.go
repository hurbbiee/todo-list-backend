package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	todoEvent "github.com/hurbbiee/todo-list-backend/internal/modules/todo/event"
)

var ErrInvalidMessage = errors.New(
	"invalid rabbitmq message",
)

type TodoCreatedNotifier interface {
	NotifyTodoCreated(
		ctx context.Context,
		event todoEvent.TodoCreated,
	) error
}

type TodoCreatedHandler struct {
	notifier TodoCreatedNotifier
}

func NewTodoCreatedHandler(
	notifier TodoCreatedNotifier,
) *TodoCreatedHandler {
	return &TodoCreatedHandler{notifier: notifier}
}

func (h *TodoCreatedHandler) Handle(
	ctx context.Context,
	body []byte,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var event todoEvent.TodoCreated

	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf(
			"%w: decode todo.created: %v",
			ErrInvalidMessage,
			err,
		)
	}

	if err := validateTodoCreated(event); err != nil {
		return fmt.Errorf(
			"%w: %v",
			ErrInvalidMessage,
			err,
		)
	}

	return h.notifier.NotifyTodoCreated(ctx, event)
}

func validateTodoCreated(
	event todoEvent.TodoCreated,
) error {
	if event.EventType != todoEvent.TodoCreatedRoutingKey {
		return fmt.Errorf(
			"unexpected eventType %q",
			event.EventType,
		)
	}

	if event.TodoID <= 0 {
		return errors.New(
			"todoId must be greater than zero",
		)
	}

	if event.UserID <= 0 {
		return errors.New(
			"userId must be greater than zero",
		)
	}

	if strings.TrimSpace(event.Title) == "" {
		return errors.New(
			"title is required",
		)
	}

	switch event.Status {
	case "pending", "in_progress", "completed":
	default:
		return fmt.Errorf(
			"unsupported status %q",
			event.Status,
		)
	}

	switch event.Priority {
	case "low", "normal", "high":
	default:
		return fmt.Errorf(
			"unsupported priority %q",
			event.Priority,
		)
	}

	if event.DueAt.IsZero() {
		return errors.New(
			"dueAt is required",
		)
	}

	if event.OccurredAt.IsZero() {
		return errors.New(
			"occurredAt is required",
		)
	}

	return nil
}
