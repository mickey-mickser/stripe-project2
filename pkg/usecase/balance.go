package usecase

import (
	"context"
	"github.com/mickey-mickser/stripe-project2"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
)

type BalanceUseCase struct {
	repo repository.UserBalance
}

func NewBalanceUseCase(repo repository.UserBalance) *BalanceUseCase {
	return &BalanceUseCase{repo: repo}
}
func (s *BalanceUseCase) GetBalance(ctx context.Context, username string) (api.User, error) {
	return s.repo.GetBalance(ctx, username)
}
func (s *BalanceUseCase) UpdateUserBalance(ctx context.Context, balance float64, username string) (float64, error) {
	return s.repo.UpdateUserBalance(ctx, balance, username)
}
