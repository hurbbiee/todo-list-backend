package repository

import (
	"context"

	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/dto"
)

type AuthRepository interface {
	FindPasswordByEmail(ctx context.Context, email string) (*dto.User, error)
}
