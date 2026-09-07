package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/model"
	todoEvent "github.com/hurbbiee/todo-list-backend/internal/modules/todo/event"
	"github.com/hurbbiee/todo-list-backend/internal/platform/discord"
)

type todoCreatedRepositoryStub struct {
	connection *model.DiscordConnection
	err        error
}

func (s *todoCreatedRepositoryStub) FindByUserID(
	_ context.Context,
	_ int64,
) (*model.DiscordConnection, error) {
	return s.connection, s.err
}

func (s *todoCreatedRepositoryStub) Upsert(
	_ context.Context,
	_ model.DiscordConnection,
) error {
	return nil
}

func (s *todoCreatedRepositoryStub) UpdateEnabled(
	_ context.Context,
	_ int64,
	_ bool,
) (bool, error) {
	return false, nil
}

func (s *todoCreatedRepositoryStub) DeleteByUserID(
	_ context.Context,
	_ int64,
) (bool, error) {
	return false, nil
}

type todoCreatedCipherStub struct {
	plaintext string
	err       error
}

func (s *todoCreatedCipherStub) Encrypt(string) ([]byte, error) {
	return nil, nil
}

func (s *todoCreatedCipherStub) Decrypt([]byte) (string, error) {
	return s.plaintext, s.err
}

type discordSenderStub struct {
	called     bool
	webhookURL string
	embed      discord.Embed
	err        error
}

func (s *discordSenderStub) SendEmbed(
	_ context.Context,
	webhookURL string,
	embed discord.Embed,
) error {
	s.called = true
	s.webhookURL = webhookURL
	s.embed = embed
	return s.err
}

func validTodoCreatedEvent() todoEvent.TodoCreated {
	return todoEvent.TodoCreated{
		EventType:  todoEvent.TodoCreatedRoutingKey,
		TodoID:     10,
		UserID:     7,
		Title:      "เขียน unit test",
		Status:     "pending",
		Priority:   "high",
		DueAt:      time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC),
		OccurredAt: time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC),
	}
}

func TestTodoCreatedNotificationServiceSendsEnabledConnection(t *testing.T) {
	repo := &todoCreatedRepositoryStub{
		connection: &model.DiscordConnection{
			UserID:              7,
			WebhookURLEncrypted: []byte("encrypted"),
			IsEnabled:           true,
			NotifyTodoCreated:   true,
		},
	}
	sender := &discordSenderStub{}
	svc := NewTodoCreatedNotificationService(
		repo,
		&todoCreatedCipherStub{plaintext: "https://discord.com/api/webhooks/id/token"},
		sender,
	)

	if err := svc.NotifyTodoCreated(context.Background(), validTodoCreatedEvent()); err != nil {
		t.Fatalf("NotifyTodoCreated() error = %v", err)
	}
	if !sender.called {
		t.Fatal("NotifyTodoCreated() did not call Discord sender")
	}
	if sender.webhookURL != "https://discord.com/api/webhooks/id/token" {
		t.Fatal("NotifyTodoCreated() sent to an unexpected webhook URL")
	}
	if !strings.Contains(sender.embed.Description, "เขียน unit test") {
		t.Fatalf("Discord embed description = %q, want todo title", sender.embed.Description)
	}
	if len(sender.embed.Fields) != 3 {
		t.Fatalf("Discord embed fields = %d, want 3", len(sender.embed.Fields))
	}
	if sender.embed.Fields[2].Value != "2026-09-08 09:00:00" {
		t.Fatalf("Discord due date = %q, want formatted date", sender.embed.Fields[2].Value)
	}
	if sender.embed.Color != discordColorRed {
		t.Fatalf("Discord embed color = %d, want high-priority red", sender.embed.Color)
	}
}

func TestTodoCreatedNotificationServiceSkipsDisabledConnection(t *testing.T) {
	sender := &discordSenderStub{}
	svc := NewTodoCreatedNotificationService(
		&todoCreatedRepositoryStub{
			connection: &model.DiscordConnection{
				IsEnabled:         false,
				NotifyTodoCreated: true,
			},
		},
		&todoCreatedCipherStub{},
		sender,
	)

	if err := svc.NotifyTodoCreated(context.Background(), validTodoCreatedEvent()); err != nil {
		t.Fatalf("NotifyTodoCreated() error = %v", err)
	}
	if sender.called {
		t.Fatal("NotifyTodoCreated() sent a notification for a disabled connection")
	}
}

func TestTodoCreatedNotificationServiceMarksDecryptFailurePermanent(t *testing.T) {
	svc := NewTodoCreatedNotificationService(
		&todoCreatedRepositoryStub{
			connection: &model.DiscordConnection{
				IsEnabled:           true,
				NotifyTodoCreated:   true,
				WebhookURLEncrypted: []byte("invalid"),
			},
		},
		&todoCreatedCipherStub{err: errors.New("decrypt failed")},
		&discordSenderStub{},
	)

	err := svc.NotifyTodoCreated(context.Background(), validTodoCreatedEvent())
	if !errors.Is(err, ErrPermanentNotification) {
		t.Fatalf("NotifyTodoCreated() error = %v, want permanent error", err)
	}
}
