package usecase

import (
	"context"
	"github.com/mickey-mickser/stripe-project2"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
)

type UserUseCase struct {
	repo repository.User
}

func NewAuthService(repo repository.User) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (s *UserUseCase) CreateUser(ctx context.Context, user api.User) (int, error) {

	return s.repo.CreateUser(ctx, user)
}
func (s *UserUseCase) GetUser(ctx context.Context, username string) (api.User, error) {
	return s.repo.GetUser(ctx, username)
}
