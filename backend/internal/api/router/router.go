package router

import (
	"net/http"
	"log/slog"

	"github.com/gin-gonic/gin"

	"qrgen/internal/api/handler"
	"qrgen/internal/api/middleware"
	"qrgen/internal/consts"
)

func NewRouter(h *handler.Handler, logger *slog.Logger, timeout int) *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.TraceMiddleware(),
		middleware.LoggingMiddleware(logger),
		middleware.RecoverMiddleware(logger),
		middleware.TimeoutMiddleware(timeout),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/gen", h.CreateQRHandler)
	r.GET("/code/:"+string(consts.ID), h.GetQRHandler)

	
	return r
}
