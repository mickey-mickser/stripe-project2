package run

import (
	"context"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"time"
)

type Syncer struct {
	Repo repository.Repository
}

func (s *Syncer) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 40)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Infof("Syncer stopped: %v", ctx.Err())

			return
		case <-ticker.C:
			sessions, err := s.Repo.SelectSessions(ctx, "open")
			if err != nil {
				logrus.Errorf("Error selecting sessions: %v", err)
				continue
			}
			if sessions == nil {
				logrus.Warn("No sessions found")
				continue
			}

			for _, val := range sessions {
				ses, err := session.Get(val.SessionID, nil)
				if err != nil {
					logrus.Errorf("Failed to get session from Stripe: %v", err)
					time.Sleep(time.Second)
					continue
				}

				if err := s.Repo.UpdateSessionStatus(ctx, val.SessionID, string(ses.Status)); err != nil {
					logrus.Errorf("Failed to update session status in db: %v", err)
					continue
				}

				switch ses.Status {
				case stripe.CheckoutSessionStatusComplete:

					userBalance, err := s.Repo.GetBalance(ctx, val.Username)
					if err != nil {
						logrus.Errorf("Failed to get user balance: %v", err)
						continue
					}
					balance := float64(val.Amount) + userBalance.Balance
					if _, err := s.Repo.UpdateUserBalance(ctx, balance, val.Username); err != nil {
						logrus.Errorf("Failed to update user balance: %v", err)
						continue
					}
					logrus.Infof("User %s balance updated successfully", val.Username)
				case stripe.CheckoutSessionStatusExpired:
					logrus.Warnf("Session %s has expired", val.SessionID)

				case stripe.CheckoutSessionStatusOpen:
					logrus.Infof("Session %s is still open", val.SessionID)
				}
			}
		}
	}
}
