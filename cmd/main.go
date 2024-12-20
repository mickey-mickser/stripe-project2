package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/mickey-mickser/stripe-project2"
	"github.com/mickey-mickser/stripe-project2/pkg/handler"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
	"github.com/mickey-mickser/stripe-project2/pkg/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error write .env: %s", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to init db: %v", err.Error())
	}

	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	//Dependency Injection(Внедрение зависимостей)
	repos := repository.NewRepository(db)
	services := usecase.NewUseCase(repos)
	//sessionStorage := cash.NewSessionStorage()
	handlers := handler.NewHandler(services)
	//handlers := handler.NewHandler(services, sessionStorage)

	// Инициализация Handler с SessionStorage
	srv := new(api.Server)

	if err := srv.Start(viper.GetString("port"), handlers.InitRouter()); err != nil {
		logrus.Fatalf("error occurred while running the HTTP server: %s", err.Error())
	}

}
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
