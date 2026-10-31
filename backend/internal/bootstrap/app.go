package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"qrgen/internal/api/handler"
	"qrgen/internal/api/router"
	"qrgen/internal/config"
	"qrgen/internal/logger"
	"qrgen/internal/repository/redisrepo"
	"qrgen/internal/service"
)

type App struct {
	AppLogger *slog.Logger
	Router    *gin.Engine
	Config    *config.Config
	Client    *redis.Client
}

func InitApp(parent context.Context) (*App, error) {
	appLogger := logger.NewLogger()

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	configEnv := ""
	switch os.Getenv("ENV") {
	case "local":
		configEnv = "CONFIG_PATH_BACKEND"
	case "dev":
		configEnv = "CONFIG_PATH_DOCKER"
	default:
		return nil, fmt.Errorf("unsupported ENV %q (expected local or dev)", os.Getenv("ENV"))
	}

	configPath := os.Getenv(configEnv)
	if configPath == "" {
		return nil, fmt.Errorf("%s is not set", configEnv)
	}

	cfg, err := config.InitConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("initialize config: %w", err)
	}

	connectCtx, cancel := context.WithTimeout(parent, time.Duration(cfg.Database.Timeout)*time.Second)
	defer cancel()

	rdb, err := redisrepo.InitClient(connectCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize Redis: %w", err)
	}

	repo := redisrepo.NewRepo(rdb)
	qrService := service.NewService(repo, time.Duration(cfg.Database.TTL)*time.Hour)
	qrHandler := handler.NewHandler(qrService)

	app := &App{
		AppLogger: appLogger,
		Router:    router.NewRouter(qrHandler),
		Config:    &cfg,
		Client:    rdb,
	}

	return app, nil
}

func (a *App) Close() error {
	if a.Client == nil {
		return nil
	}
	return a.Client.Close()
}
