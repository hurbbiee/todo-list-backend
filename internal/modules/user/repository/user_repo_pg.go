package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hurbbiee/todo-list-backend/internal/modules/user/dto"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoPg struct {
	db *pgxpool.Pool
}

func NewUserRepoPg(db *pgxpool.Pool) UserRepository {
	return &UserRepoPg{db: db}
}

func (r *UserRepoPg) Create(
	ctx context.Context,
	req dto.CreateUserRequest,
	password string,
) error {

	creatQuery := `
					INSERT INTO "users"
						("name",email,password_hash,"role")
					VALUES 
						($1,$2,$3,$4)`

	cmdTag, err := r.db.Exec(
		ctx,
		creatQuery,
		req.Name,
		req.Email,
		password,
		req.Role,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return enum.ErrEmailAlreadyExists
			}
		}
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user can't create")
	}

	return nil
}

func (r *UserRepoPg) GetProfile(
	ctx context.Context,
	id int64,
) (dto.ProfileResponse, error) {

	var user dto.ProfileResponse

	baseQuery := `
				SELECT 
					u.id,
					u."name",
					u.email,
					u.role
				FROM "users" u
				WHERE u.id = $1
				AND u.is_deleted = false`

	err := r.db.QueryRow(ctx, baseQuery, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}

	return user, nil
}
