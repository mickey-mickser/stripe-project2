package repository

import (
	"context"
	"errors"
	api "github.com/mickey-mickser/stripe-project2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type SessionStatus struct {
	db *gorm.DB
}

func NewtSessionStatus(db *gorm.DB) *SessionStatus {
	return &SessionStatus{db: db}
}

type PaymentSession struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	SessionID string    `gorm:"type:varchar(255);not null" json:"session_id"`
	Username  string    `gorm:"type:varchar(255);not null" json:"username"`
	Amount    float64   `gorm:"type:numeric(10,2);not null" json:"amount"`
	Status    string    `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (s *SessionStatus) CreateSession(ctx context.Context, sessionID, username, status string, amount float64) error {
	session := PaymentSession{
		SessionID: sessionID,
		Username:  username,
		Amount:    amount,
		Status:    status,
	}

	if err := s.db.WithContext(ctx).Create(&session).Error; err != nil {
		return err
	}

	return nil
}

func (s *SessionStatus) UpdateSessionStatus(ctx context.Context, sessionID, status string) error {
	if err := s.db.WithContext(ctx).
		Model(&PaymentSession{}).
		Where("session_id = ?", sessionID).
		Update("status", status).Error; err != nil {
		logrus.Errorf("Failed to update session status in DB: %v", err)
		return err
	}

	logrus.Infof("Session status for session ID %s updated to %s", sessionID, status)
	return nil
}

func (s *SessionStatus) GetStatus(ctx context.Context, sessionID string) (*api.PaymentSession, error) {
	var status api.PaymentSession

	query := s.db.WithContext(ctx).Model(&api.PaymentSession{}).
		Where("session_id = ?", sessionID).
		First(&status)

	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}

		return nil, err
	}

	return &status, nil
}
func (s *SessionStatus) SelectSessions(ctx context.Context, status string) ([]api.PaymentSession, error) {
	var sessions []api.PaymentSession
	query := s.db.WithContext(ctx).Model(&api.PaymentSession{}).
		Where("status = ?", status).
		Find(&sessions)
	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return sessions, nil
}
