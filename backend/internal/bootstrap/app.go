package bootstrap

import "qrgen/internal/logger"


func InitApp() {
	logger := logger.NewLogger()
	_ = logger
}