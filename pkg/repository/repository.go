package repository

import (
	"context"
	"github.com/mickey-mickser/stripe-project2"
	"gorm.io/gorm"
)

type Session interface {
	CreateSession(ctx context.Context, sessionID, username, status string, amount float64) error
	UpdateSessionStatus(ctx context.Context, sessionID, status string) error
	GetStatus(ctx context.Context, sessionID string) (*api.PaymentSession, error)
	SelectSessions(ctx context.Context, status string) ([]api.PaymentSession, error)
}
type UserBalance interface {
	GetBalance(ctx context.Context, username string) (api.User, error)
	UpdateUserBalance(ctx context.Context, balance float64, username string) (float64, error)
}
type User interface {
	CreateUser(ctx context.Context, user api.User) (int, error)
	GetUser(ctx context.Context, username string) (api.User, error)
}

type Repository struct {
	User
	UserBalance
	Session
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:        NewAuthPostgres(db),
		UserBalance: NewBalancePostgres(db),
		Session:     NewtSessionStatus(db),
	}
}
