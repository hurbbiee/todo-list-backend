package repository

import (
	"context"

	"github.com/hurbbiee/todo-list-backend/internal/modules/notification/model"
)

type DiscordConnectionRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*model.DiscordConnection, error)
	Upsert(ctx context.Context, connection model.DiscordConnection) error
	UpdateEnabled(ctx context.Context, userID int64, isEnabled bool) (bool, error)
	DeleteByUserID(ctx context.Context, userID int64) (bool, error)
}
