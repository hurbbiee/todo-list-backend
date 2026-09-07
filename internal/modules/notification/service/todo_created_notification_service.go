package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/repository"
	todoEvent "github.com/hurbbiee/todo-list-backend/internal/modules/todo/event"
	"github.com/hurbbiee/todo-list-backend/internal/platform/discord"
)

const (
	maxDiscordTodoTitleRunes = 1500
	discordColorGreen        = 0x57F287
	discordColorYellow       = 0xFEE75C
	discordColorRed          = 0xED4245
)

var ErrPermanentNotification = errors.New("permanent notification error")

type DiscordSender interface {
	SendEmbed(
		ctx context.Context,
		webhookURL string,
		embed discord.Embed,
	) error
}

type TodoCreatedNotificationService struct {
	repo   repository.DiscordConnectionRepository
	cipher SecretCipher
	sender DiscordSender
}

func NewTodoCreatedNotificationService(
	repo repository.DiscordConnectionRepository,
	cipher SecretCipher,
	sender DiscordSender,
) *TodoCreatedNotificationService {
	return &TodoCreatedNotificationService{
		repo:   repo,
		cipher: cipher,
		sender: sender,
	}
}

func (s *TodoCreatedNotificationService) NotifyTodoCreated(
	ctx context.Context,
	event todoEvent.TodoCreated,
) error {
	connection, err := s.repo.FindByUserID(ctx, event.UserID)
	if err != nil {
		return fmt.Errorf("find Discord connection: %w", err)
	}

	if connection == nil ||
		!connection.IsEnabled ||
		!connection.NotifyTodoCreated {
		return nil
	}

	webhookURL, err := s.cipher.Decrypt(connection.WebhookURLEncrypted)
	if err != nil {
		return fmt.Errorf(
			"%w: decrypt Discord webhook URL: %v",
			ErrPermanentNotification,
			err,
		)
	}

	embed := formatTodoCreatedEmbed(event)
	if err := s.sender.SendEmbed(ctx, webhookURL, embed); err != nil {
		return fmt.Errorf("send Discord todo.created notification: %w", err)
	}

	return nil
}

func formatTodoCreatedEmbed(event todoEvent.TodoCreated) discord.Embed {
	title := truncateRunes(event.Title, maxDiscordTodoTitleRunes)
	dueAt := event.DueAt.Format("2006-01-02 15:04:05")

	return discord.Embed{
		Title:       "📌 สร้าง Todo ใหม่",
		Description: title,
		Color:       todoPriorityColor(event.Priority),
		Fields: []discord.EmbedField{
			{Name: "📊 สถานะ", Value: event.Status, Inline: true},
			{Name: "🔥 ความสำคัญ", Value: event.Priority, Inline: true},
			{Name: "⏰ กำหนดส่ง", Value: dueAt, Inline: false},
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Todo ID: %d", event.TodoID),
		},
	}
}

func todoPriorityColor(priority string) int {
	switch priority {
	case "high":
		return discordColorRed
	case "normal":
		return discordColorYellow
	default:
		return discordColorGreen
	}
}

func truncateRunes(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}

	return string(runes[:limit]) + "…"
}
