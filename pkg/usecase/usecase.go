package usecase

import (
	"context"
	"github.com/mickey-mickser/stripe-project2"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
)

type Session interface {
	CreateSession(ctx context.Context, sessionID, username, status string, amount float64) error
	UpdateSessionStatus(ctx context.Context, sessionID, status string) error
	GetStatus(ctx context.Context, sessionID string) (*api.PaymentSession, error)
}
type UserBalance interface {
	GetBalance(ctx context.Context, username string) (api.User, error)
	UpdateUserBalance(ctx context.Context, balance float64, username string) (float64, error)
}

type User interface {
	CreateUser(ctx context.Context, user api.User) (int, error)
	GetUser(ctx context.Context, username string) (api.User, error)
}

type UseCase struct {
	User
	UserBalance
	Session
}

func NewUseCase(repo *repository.Repository) *UseCase {
	return &UseCase{
		User:        NewAuthService(repo.User),
		UserBalance: NewBalanceUseCase(repo.UserBalance),
		Session:     NewtSessionStatus(repo.Session),
	}
}
