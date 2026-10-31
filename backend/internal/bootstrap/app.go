package bootstrap

import (
	"fmt"
	"log/slog"
	"os"
	"qrgen/internal/config"
	"qrgen/internal/logger"
	"qrgen/internal/repository/postgres"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	AppLogger *slog.Logger
	Router *gin.Engine
}

func InitApp() (*App, error){
	app := &App{}
	appLogger := logger.NewLogger()	

	err := godotenv.Load()

	if err != nil {
		appLogger.Error(
			"godotenv error:",
			"error", err,
		)
		return nil, fmt.Errorf("godotenv error %v \n", err)
	}
	
	configEnv := ""
	env := os.Getenv("ENV")
	switch env {
	case "local":
		configEnv = "CONFIG_PATH_BACKEND"
	case "dev":
		configEnv = "CONFIG_PATH_DOCKER"
	default:
		appLogger.Error(
			"wrong work env",
			"error", "wrong config_env",
			"env", env,
		)
		return nil, fmt.Errorf("wrong work env %v\n", err)
	}

	path := os.Getenv(configEnv)

	cfg, err := config.InitConfig(path)
	if err != nil {
		appLogger.Error(
			"config init failed",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("config init failed %v \n", err)
	}

	db, err := postgres.InitDB(cfg)
	if err != nil {
		appLogger.Error(
			"db init failed",
			"error", err.Error(),
		)
		return nil, fmt.Errorf("db init failed %v \n", err)
	}

	_ = db
	return app, nil
}


