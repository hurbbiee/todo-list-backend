package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	todoEvent "github.com/hurbbiee/todo-list-backend/internal/modules/todo/event"
)

type todoCreatedNotifierStub struct {
	called bool
	event  todoEvent.TodoCreated
	err    error
}

func (s *todoCreatedNotifierStub) NotifyTodoCreated(
	_ context.Context,
	event todoEvent.TodoCreated,
) error {
	s.called = true
	s.event = event
	return s.err
}

func validTodoCreatedBody(t *testing.T) []byte {
	t.Helper()

	body, err := json.Marshal(todoEvent.TodoCreated{
		EventType:  todoEvent.TodoCreatedRoutingKey,
		TodoID:     10,
		UserID:     7,
		Title:      "test todo",
		Status:     "pending",
		Priority:   "normal",
		DueAt:      time.Now().Add(time.Hour),
		OccurredAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("marshal todo.created: %v", err)
	}

	return body
}

func TestTodoCreatedHandlerCallsNotifier(t *testing.T) {
	notifier := &todoCreatedNotifierStub{}
	handler := NewTodoCreatedHandler(notifier)

	if err := handler.Handle(context.Background(), validTodoCreatedBody(t)); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if !notifier.called || notifier.event.UserID != 7 {
		t.Fatal("Handle() did not forward the decoded event")
	}
}

func TestTodoCreatedHandlerRejectsInvalidMessage(t *testing.T) {
	notifier := &todoCreatedNotifierStub{}
	handler := NewTodoCreatedHandler(notifier)

	err := handler.Handle(context.Background(), []byte("not-json"))
	if !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("Handle() error = %v, want invalid message", err)
	}
	if notifier.called {
		t.Fatal("Handle() called notifier for an invalid message")
	}
}

func TestTodoCreatedHandlerReturnsNotifierError(t *testing.T) {
	wantErr := errors.New("send failed")
	handler := NewTodoCreatedHandler(&todoCreatedNotifierStub{err: wantErr})

	err := handler.Handle(context.Background(), validTodoCreatedBody(t))
	if !errors.Is(err, wantErr) {
		t.Fatalf("Handle() error = %v, want %v", err, wantErr)
	}
}
