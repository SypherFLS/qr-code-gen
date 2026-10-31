package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qrgen/internal/api/handler"
	"qrgen/internal/consts"
)

func NewRouter(h *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/gen", h.CreateQRHandler)
	r.GET("/code/:"+string(consts.ID), h.GetQRHandler)

	return r
}
