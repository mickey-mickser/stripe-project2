package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mickey-mickser/stripe-project2"
	"github.com/mickey-mickser/stripe-project2/pkg/handler"
	"github.com/mickey-mickser/stripe-project2/pkg/repository"
	"github.com/mickey-mickser/stripe-project2/pkg/usecase"
	"github.com/mickey-mickser/stripe-project2/run"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error write .env: %s", err)
	}

	db, sqlDB, err := repository.NewPostgresDB(repository.Config{
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
	defer repository.ClosePostgresDB(sqlDB)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	dataFromSession := make(chan handler.DataSession)
	logrus.Info(dataFromSession)
	//Dependency Injection(Внедрение зависимостей)
	repos := repository.NewRepository(db)
	services := usecase.NewUseCase(repos)
	handlers := handler.NewHandler(services)
	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syn := &run.Syncer{
		Repo: *repos,
	}
	go syn.Start(ctx)
	srv := new(api.Server)
	go func() {
		if err := srv.Start(viper.GetString("port"), handlers.InitRouter()); err != nil {
			logrus.Fatalf("error occurred while running the HTTP server: %s", err.Error())
		}
	}()
	logrus.Print("Api Started")
	/////////////////

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Api Down")
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occurred on server shutdown: %s", err.Error())
	}

}
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
