package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DiscordConnectionRepoPg struct {
	db *pgxpool.Pool
}

func NewDiscordConnectionRepoPg(
	db *pgxpool.Pool,
) DiscordConnectionRepository {
	return &DiscordConnectionRepoPg{db: db}
}

func (r *DiscordConnectionRepoPg) FindByUserID(
	ctx context.Context,
	userID int64,
) (*model.DiscordConnection, error) {
	const query = `
		SELECT
			id,
			user_id,
			webhook_url_encrypted,
			is_enabled,
			notify_todo_created,
			notify_todo_completed,
			notify_before_due,
			remind_before_minutes,
			created_at,
			updated_at
		FROM discord_connections
		WHERE user_id = $1
		LIMIT 1
	`

	var connection model.DiscordConnection
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&connection.ID,
		&connection.UserID,
		&connection.WebhookURLEncrypted,
		&connection.IsEnabled,
		&connection.NotifyTodoCreated,
		&connection.NotifyTodoCompleted,
		&connection.NotifyBeforeDue,
		&connection.RemindBeforeMinutes,
		&connection.CreatedAt,
		&connection.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find Discord connection by user ID: %w", err)
	}

	return &connection, nil
}

func (r *DiscordConnectionRepoPg) Upsert(
	ctx context.Context,
	connection model.DiscordConnection,
) error {
	const query = `
		INSERT INTO discord_connections (
			user_id,
			webhook_url_encrypted,
			is_enabled,
			notify_todo_created,
			notify_todo_completed,
			notify_before_due,
			remind_before_minutes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id)
		DO UPDATE SET
			webhook_url_encrypted = EXCLUDED.webhook_url_encrypted,
			is_enabled = EXCLUDED.is_enabled,
			notify_todo_created = EXCLUDED.notify_todo_created,
			notify_todo_completed = EXCLUDED.notify_todo_completed,
			notify_before_due = EXCLUDED.notify_before_due,
			remind_before_minutes = EXCLUDED.remind_before_minutes,
			updated_at = NOW()
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		connection.UserID,
		connection.WebhookURLEncrypted,
		connection.IsEnabled,
		connection.NotifyTodoCreated,
		connection.NotifyTodoCompleted,
		connection.NotifyBeforeDue,
		connection.RemindBeforeMinutes,
	)
	if err != nil {
		return fmt.Errorf("upsert Discord connection: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf(
			"upsert Discord connection: expected 1 affected row, got %d",
			commandTag.RowsAffected(),
		)
	}

	return nil
}

func (r *DiscordConnectionRepoPg) UpdateEnabled(
	ctx context.Context,
	userID int64,
	isEnabled bool,
) (bool, error) {
	const query = `
		UPDATE discord_connections
		SET
			is_enabled = $2,
			updated_at = NOW()
		WHERE user_id = $1
	`

	commandTag, err := r.db.Exec(ctx, query, userID, isEnabled)
	if err != nil {
		return false, fmt.Errorf("update Discord connection status: %w", err)
	}

	return commandTag.RowsAffected() == 1, nil
}

func (r *DiscordConnectionRepoPg) DeleteByUserID(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const query = `
		DELETE FROM discord_connections
		WHERE user_id = $1
	`

	commandTag, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return false, fmt.Errorf("delete Discord connection: %w", err)
	}

	return commandTag.RowsAffected() == 1, nil
}
