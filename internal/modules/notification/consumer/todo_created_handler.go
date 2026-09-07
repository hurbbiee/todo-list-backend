package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	todoEvent "github.com/hurbbiee/todo-list-backend/internal/modules/todo/event"
)

var ErrInvalidMessage = errors.New(
	"invalid rabbitmq message",
)

type TodoCreatedHandler struct{}

func NewTodoCreatedHandler() *TodoCreatedHandler {
	return &TodoCreatedHandler{}
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

	// ตอนนี้ log ข้อมูลที่ parse แล้ว
	// ขั้นต่อไปตรงนี้จะเปลี่ยนเป็นส่ง Discord
	log.Printf(
		"todo.created handled: todoId=%d userId=%d title=%q status=%s priority=%s dueAt=%s",
		event.TodoID,
		event.UserID,
		event.Title,
		event.Status,
		event.Priority,
		event.DueAt.Format("2006-01-02 15:04:05Z07:00"),
	)

	return nil
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