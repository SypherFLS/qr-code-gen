package router

import (
	"qrgen/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handler) *gin.Engine{
	r := gin.New()

	r.POST("/gen", h.CreateHandler)
	r.GET("/code", h.GetHandler)
	
	return r
}