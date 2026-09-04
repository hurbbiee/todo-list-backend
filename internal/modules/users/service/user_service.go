package service

import (
	"context"

	"github.com/hurbbiee/todo-list-backend/internal/modules/users/dto"
	"github.com/hurbbiee/todo-list-backend/internal/modules/users/repository"
	bcrypt "github.com/hurbbiee/todo-list-backend/internal/platform/security"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(
	repo repository.UserRepository,
) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) error {
	passwordHash, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, req, passwordHash)
}

func (s *UserService) Profile(
	ctx context.Context,
	id int64,
) (dto.ProfileResponse, error) {

	user, err := s.repo.GetProfile(ctx, id)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return user, nil
}
