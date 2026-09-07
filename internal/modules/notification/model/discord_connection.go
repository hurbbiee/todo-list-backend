package model

import "time"

type DiscordConnection struct {
	ID                  int64
	UserID              int64
	WebhookURLEncrypted []byte
	IsEnabled           bool
	NotifyTodoCreated   bool
	NotifyTodoCompleted bool
	NotifyBeforeDue     bool
	RemindBeforeMinutes int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
