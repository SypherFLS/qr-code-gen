package middleware

import (
	"context"
	// "log"
	"log/slog"
	"net/http"
	// "qrgen/internal/api/utils"
	"qrgen/internal/consts"
	"time"
	"uuid"

	"github.com/gin-gonic/gin"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		ctx := context.WithValue(
			c.Request.Context(),
			consts.RequestIDKey,
			requestID,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// requestID := utils.GetRequestID(c.Request.Context())
		c.Next()
		
		latency := time.Since(start)

		logger.Info("Request took %v", latency)
	}
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(){
			if err := recover(); err != nil {
				// logger.Error("", )

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		

		c.Next()
	}
}

