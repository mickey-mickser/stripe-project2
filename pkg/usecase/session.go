package usecase

import (
	"context"
	api "github.com/mickey-mickser/stripe-project2"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
)

type SessionStatus struct {
	repo repository.Session
}

func NewtSessionStatus(repo repository.Session) *SessionStatus {
	return &SessionStatus{repo: repo}
}
func (s *SessionStatus) CreateSession(ctx context.Context, sessionID, username, status string, amount float64) error {
	return s.repo.CreateSession(ctx, sessionID, username, status, amount)
}
func (s *SessionStatus) UpdateSessionStatus(ctx context.Context, sessionID, status string) error {
	return s.repo.UpdateSessionStatus(ctx, sessionID, status)
}

func (s *SessionStatus) GetStatus(ctx context.Context, sessionID string) (*api.PaymentSession, error) {
	return s.repo.GetStatus(ctx, sessionID)
}
