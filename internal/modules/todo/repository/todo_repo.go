package repository

import (
	"context"

	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/dto"
)

type TodoRepository interface {
	Search(ctx context.Context, req dto.SearchTodoRequest, actionBy int64) ([]dto.SearchTodoResponse, int64, error)
	Create(ctx context.Context, req dto.CreateTodoRequest, actionBy int64) (int64, error)
	CountStatus(ctx context.Context, actionBy int64) (dto.CountTodoResponse, error)
	Update(ctx context.Context, req dto.UpdateTodoRequest, id int64, actionBy int64) error
	UpdateStatus(ctx context.Context, req dto.UpdateTodoStatus, id int64, actionBy int64) error
	Delete(ctx context.Context, id int64, actionBy int64) error
}
