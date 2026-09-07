package repository

import (
	"context"
	"fmt"

	"github.com/hurbbiee/todo-list-backend/internal/modules/todo/dto"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepoPg struct {
	db *pgxpool.Pool
}

func NewTodoRepopg(db *pgxpool.Pool) TodoRepository {
	return &TodoRepoPg{db: db}
}

func (r *TodoRepoPg) Search(
	ctx context.Context,
	req dto.SearchTodoRequest,
	actionBy int64,
) ([]dto.SearchTodoResponse, int64, error) {

	var page int64 = 1
	var limit int64 = 10

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := (page - 1) * limit

	baseQuery := `
				FROM "todos" t 
				WHERE 1=1 AND t.is_deleted = false
				AND t.user_id = $1`

	args := []interface{}{actionBy}
	argID := 2

	if req.Keyword != "" {
		baseQuery += fmt.Sprintf(`
				AND (
					t.title LIKE $%d
					OR t.description LIKE $%d
				)`, argID, argID)
		args = append(args, "%"+req.Keyword+"")
		argID++
	}

	if req.Priority != nil && *req.Priority != "" && *req.Priority != "all" {
		baseQuery += fmt.Sprintf(`
				AND (
					t.priority = $%d
				)`, argID)
		args = append(args, *req.Priority)
		argID++
	}

	if req.Status != nil && *req.Status != "" && *req.Status != "all" {
		baseQuery += fmt.Sprintf(`
				AND (
					t.status = $%d
				)`, argID)
		args = append(args, *req.Status)
		argID++
	}

	if req.DateFrom != nil && !req.DateFrom.IsZero() {
		baseQuery += fmt.Sprintf(" AND t.due_at::date >= $%d", argID)
		args = append(args, req.DateFrom.Time)
		argID++
	}

	if req.DateTo != nil && !req.DateTo.IsZero() {
		baseQuery += fmt.Sprintf(" AND t.due_at::date <= $%d", argID)
		args = append(args, req.DateTo.Time)
		argID++
	}

	countQuery := `SELECT COUNT(DISTINCT t.id) ` + baseQuery
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []dto.SearchTodoResponse{}, 0, nil
	}

	selectQuery := `
				SELECT 
					t.id,
					t.title,
					t.description,
					t.priority,
					t.status,
					t.due_at,
					t.completed_at,
					t.created_via
				` + baseQuery + `
				GROUP BY t.id
				ORDER BY t.created_at DESC`

	if !req.FetchAll {
		selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []dto.SearchTodoResponse
	for rows.Next() {
		var t dto.SearchTodoResponse

		err := rows.Scan(
			&t.Id,
			&t.Title,
			&t.Description,
			&t.Priority,
			&t.Status,
			&t.DueDate,
			&t.CompletedAt,
			&t.CreatedVia,
		)
		if err != nil {
			return nil, 0, err
		}

		todos = append(todos, t)
	}

	return todos, total, nil
}

func (r *TodoRepoPg) Create(
	ctx context.Context,
	req dto.CreateTodoRequest,
	actionBy int64,
) (int64, error) {

	baseQuery := `
				INSERT INTO todos
					(user_id, title, description, status, priority, due_at, created_via,completed_at)
				VALUES
					($1,$2,$3,$4::text,$5,$6,$7,CASE
											WHEN $4::text = 'completed' THEN NOW()
											ELSE NULL
										END
				)
				RETURNING id
			`

	var todoID int64

	err := r.db.QueryRow(
		ctx,
		baseQuery,
		actionBy,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		req.DueDate,
		req.CreatedVia,
	).Scan(&todoID)

	if err != nil {
		return 0, fmt.Errorf(
			"create todo: %w",
			err,
		)
	}

	return todoID, nil
}

func (r *TodoRepoPg) CountStatus(
	ctx context.Context,
	actionBy int64,
) (dto.CountTodoResponse, error) {
	query := `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE t.status = 'pending'),
			COUNT(*) FILTER (WHERE t.status = 'in_progress'),
			COUNT(*) FILTER (WHERE t.status = 'completed')
		FROM todos t
		WHERE t.user_id = $1
			AND t.is_deleted = false
	`

	var result dto.CountTodoResponse

	err := r.db.QueryRow(ctx, query, actionBy).Scan(
		&result.All,
		&result.Pending,
		&result.InProgress,
		&result.Completed,
	)
	if err != nil {
		return dto.CountTodoResponse{}, err
	}

	return result, nil
}

func (r *TodoRepoPg) Update(
	ctx context.Context,
	req dto.UpdateTodoRequest,
	id int64,
	actionBy int64,
) error {

	baseQuery := `
				UPDATE todos
				SET 
					title = $3,
					description = $4,
					status = $5::text,
					priority = $6,
					updated_via = $7,
					updated_at = NOW(),
					due_at = $8,
					completed_at = CASE
						WHEN $5::text = 'completed'
							THEN COALESCE(completed_at, NOW())
						ELSE NULL
					END
				WHERE id = $1
				AND is_deleted = false
				AND user_id = $2
				`
	cmdTag, err := r.db.Exec(
		ctx,
		baseQuery,
		id,
		actionBy,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		req.UpdatedVia,
		req.DueDate,
	)
	if err != nil {
		return fmt.Errorf("update todo: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

func (r *TodoRepoPg) UpdateStatus(
	ctx context.Context,
	req dto.UpdateTodoStatus,
	id int64,
	actionBy int64,
) error {

	baseQuery := `
	UPDATE todos
	SET
		status = $3::text,
		completed_at = CASE
			WHEN $3::text = 'completed'
				THEN COALESCE(completed_at, NOW())
			ELSE NULL
		END,
		updated_via = $4,
		updated_at = NOW()
	WHERE id = $1
		AND user_id = $2
		AND is_deleted = FALSE
`

	cmdTag, err := r.db.Exec(
		ctx,
		baseQuery,
		id,
		actionBy,
		req.Status,
		req.UpdatedVia,
	)
	if err != nil {
		return fmt.Errorf("update todo status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

func (r *TodoRepoPg) Delete(
	ctx context.Context,
	id int64,
	actionBy int64,
) error {

	baseQuery := `
				UPDATE todos
				SET 
					is_deleted = true,
					deleted_at = NOW(),
					updated_at = NOW()
				WHERE id = $1
				AND is_deleted = false
				AND user_id = $2
	`

	cmdTag, err := r.db.Exec(
		ctx,
		baseQuery,
		id,
		actionBy,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return enum.ErrUserNotFound
	}

	return nil
}
