package dto

import (
	"github.com/hurbbiee/todo-list-backend/internal/shared/types"
)

type SearchTodoRequest struct {
	Keyword  string       `json:"keyword" validate:"omitempty,min=1"`
	Status   *string      `json:"status" validate:"omitempty,oneof='pending' 'in_progress' 'completed' 'all'"`
	Priority *string      `json:"priority" validate:"omitempty,oneof=high normal low all"`
	FetchAll bool         `json:"fetchAll"`
	DateFrom *types.Date  `json:"dateFrom" validate:"omitempty"`
	DateTo   *types.Date  `json:"dateTo" validate:"omitempty"`
	Page     int64        `json:"page" validate:"omitempty,min=1"`
	Sort     []SortOption `json:"sort" validate:"omitempty,dive"`
	Limit    int64        `json:"limit" validate:"omitempty,min=1,max=100"`
}

type SortOption struct {
	Field string `json:"field" validate:"required,oneof=created_at"`
	Order string `json:"order" validate:"required,oneof=asc desc"`
}

func (r *SearchTodoRequest) ApplyDefault() {
	if r.Page <= 0 {
		r.Page = 1
	}

	if r.Limit <= 0 {
		r.Limit = 10
	}

	if len(r.Sort) == 0 {
		r.Sort = []SortOption{
			{
				Field: "created_at",
				Order: "desc",
			},
		}
	}
}
