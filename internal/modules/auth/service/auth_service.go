package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hurbbiee/todo-list-backend/internal/modules/auth/repository"
	dto "github.com/hurbbiee/todo-list-backend/internal/shared/auth"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
	"github.com/hurbbiee/todo-list-backend/internal/shared/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(
	repo repository.AuthRepository,
	jwtSecret []byte,
) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.repo.FindPasswordByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", response.BadRequest(enum.AuthInvalidCredential)
	}

	claims := dto.CustomClaims{
		UserID: user.ID,
		Role:   user.Role,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "todolist-auth-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", response.InternalError(enum.Internal)
	}

	return signed, nil
}
