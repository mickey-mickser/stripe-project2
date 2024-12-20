package repository

import (
	"context"
	"errors"
	"github.com/mickey-mickser/stripe-project2"
	"gorm.io/gorm"
)

type BalancePostgres struct {
	db *gorm.DB
}

func NewBalancePostgres(db *gorm.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func (r *BalancePostgres) GetBalance(ctx context.Context, username string) (api.User, error) {
	var user api.User

	query := r.db.WithContext(ctx).Model(&user).
		Select("balance").
		Where("username = ?", username).
		First(&user)
	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.User{}, errors.New("user not found")
		}
		return api.User{}, err
	}
	return user, nil
} //query := r.db.WithContext(ctx).Model(&user).
// Select("balance").
// Where("username = ?", username).
// First(&user)
//
//	if err := query.Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return api.User{}, errors.New("user not found")
//		}
//		return api.User{}, err
//	}
//func (r *BalancePostgres) UpdateUserBalance(ctx context.Context, balance float64, username string) (float64, error) {
//	var user api.User
//
//	query := r.db.WithContext(ctx).
//		Model(&user).
//		Update("balance", balance).
//		Where("username = ?", username)
//	if err := query.Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return 0.0, errors.New("user not found")
//		}
//
//		return 0.0, errors.New("not implemented")
//	}
//	return user.Balance, nil
//}

func (r *BalancePostgres) UpdateUserBalance(ctx context.Context, balance float64, username string) (float64, error) {
	var user api.User

	query := r.db.WithContext(ctx).
		Model(&user).
		Where("username = ?", username).
		Update("balance", balance)

	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0.0, errors.New("user not found")
		}
		return 0.0, err
	}

	if err := r.db.WithContext(ctx).Model(&user).Where("username = ?", username).First(&user).Error; err != nil {
		return 0.0, err
	}

	return user.Balance, nil
}
