package repository

import (
	"context"
	"errors"

	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/dto"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepoPg struct {
	db *pgxpool.Pool
}

func NewAuthRepoPg(db *pgxpool.Pool) AuthRepository {
	return &AuthRepoPg{db: db}
}

func (r *AuthRepoPg) FindPasswordByEmail(ctx context.Context, email string) (*dto.User, error) {
	var user dto.User

	baseQuery := `
				SELECT 
					id,
					"name",
					email,
					password_hash,
					role
				FROM "users"
				WHERE email = $1
				AND is_deleted = false
				LIMIT 1`

	err := r.db.QueryRow(
		ctx,
		baseQuery,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, response.BadRequest(enum.AuthInvalidCredential)
	}

	if err != nil {
		return nil, response.InternalError(enum.Internal)
	}

	return &user, nil
}
