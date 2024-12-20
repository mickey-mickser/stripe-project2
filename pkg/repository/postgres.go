package repository

import (
	"database/sql"
	"fmt"
	migrations "github.com/mickey-mickser/stripe-project2/pkg/migration"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*gorm.DB, *sql.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, nil, err
	}
	// Check migrations
	var userTableExists bool
	err = db.Table("information_schema.tables").
		Where("table_name = ?", "users").
		Where("table_schema = 'public'").
		Select("1").
		Limit(1).
		Scan(&userTableExists).Error

	if err != nil {
		logrus.Errorf("Error checking 'users' table existence: %v", err)
		return nil, nil, err
	}

	var paymentSessionTableExists bool
	err = db.Table("information_schema.tables").
		Where("table_name = ?", "payment_sessions").
		Where("table_schema = 'public'").
		Select("1").
		Limit(1).
		Scan(&paymentSessionTableExists).Error

	if err != nil {
		logrus.Errorf("Error checking 'payment_sessions' table existence: %v", err)
		return nil, nil, err

	}

	if !userTableExists || !paymentSessionTableExists {
		logrus.Info("One or more tables do not exist, applying migrations.")
		if err = migrations.UpMigrations(db); err != nil {
			logrus.Errorf("Error applying migrations: %v", err)
			return nil, nil, fmt.Errorf("migration application failed: %v", err)
		}
		logrus.Info("Migrations applied successfully.")
	} else {
		logrus.Info("Both tables 'users' and 'payment_sessions' already exist, skipping migrations.")
	}
	//
	return db, sqlDB, nil

}
func ClosePostgresDB(sqlDB *sql.DB) {
	if sqlDB != nil {
		sqlDB.Close()
	}
}
