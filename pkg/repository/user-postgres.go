package repository

import (
	"context"
	"errors"
	"github.com/mickey-mickser/stripe-project2"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) CreateUser(ctx context.Context, user api.User) (int, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, errors.New("operation canceled or timed out")
		}

		return 0, err
	}

	return user.Id, nil
}

func (r *AuthPostgres) GetUser(ctx context.Context, username string) (api.User, error) {
	var user api.User

	query := r.db.WithContext(ctx).Model(&api.User{}).
		Where("username = ?", username).
		First(&user)

	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.User{}, errors.New("user not found")
		}
		return api.User{}, err
	}

	return user, nil
}
