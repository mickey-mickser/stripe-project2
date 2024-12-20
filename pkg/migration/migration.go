package migrations

import (
	"github.com/mickey-mickser/stripe-project2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UpMigrations выполняет ап-миграции (создание таблиц)
func UpMigrations(db *gorm.DB) error {
	// Применяем миграцию для создания таблицы 'users'
	if err := db.AutoMigrate(&api.User{}); err != nil {
		logrus.Errorf("Migration failed for 'users' table: %v", err)
		return err
	}
	logrus.Info("Table 'users' created or already exists.")

	// Применяем миграцию для создания таблицы 'payment_sessions'
	if err := db.AutoMigrate(&api.PaymentSession{}); err != nil {
		logrus.Errorf("Migration failed for 'payment_sessions' table: %v", err)
		return err
	}
	logrus.Info("Table 'payment_sessions' created or already exists.")

	return nil
}

// DownMigrations выполняет даун-миграции (удаление таблиц)
func DownMigrations(db *gorm.DB) error {
	// Удаляем таблицу 'users'
	if err := db.Migrator().DropTable(&api.User{}); err != nil {
		logrus.Errorf("Failed to drop 'users' table: %v", err)
		return err
	}
	logrus.Info("Table 'users' dropped.")

	// Удаляем таблицу 'payment_sessions'
	if err := db.Migrator().DropTable(&api.PaymentSession{}); err != nil {
		logrus.Errorf("Failed to drop 'payment_sessions' table: %v", err)
		return err
	}
	logrus.Info("Table 'payment_sessions' dropped.")

	return nil
}
