package repository

import (
	"context"

	"github.com/hurbbiee/todo-list-backend/internal/modules/users/dto"
)

type UserRepository interface {
	Create(ctx context.Context, req dto.CreateUserRequest, password string) error
	GetProfile(ctx context.Context, id int64) (dto.ProfileResponse, error)
}
