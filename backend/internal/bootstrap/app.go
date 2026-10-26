package bootstrap

import "qrgen/internal/logger"


func Init() {
	logger := logger.NewLogger()
	_ = logger
}