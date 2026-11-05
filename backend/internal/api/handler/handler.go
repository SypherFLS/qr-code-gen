package handler

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"qrgen/internal/apperrors"
	"qrgen/internal/consts"
	"qrgen/internal/dto"
	"qrgen/internal/application"
)

type Handler struct {
	application *application.Application
}

func NewHandler(application *application.Application) *Handler {
	return &Handler{application: application}
}

func (h *Handler) CreateQRHandler(c *gin.Context) {
	var req dto.CreateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "url must be a valid HTTP or HTTPS URL"})
		return
	}

	id, png, err := h.application.CreateQRApplication(c.Request.Context(), req.URL)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create QR code"})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateQRResponse{
		ID:    id,
		Image: base64.StdEncoding.EncodeToString(png),
	})
}

func (h *Handler) GetQRHandler(c *gin.Context) {
	id := c.Param(string(consts.ID))
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id is required"})
		return
	}

	png, err := h.application.GetQRApplication(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "QR code not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to retrieve QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}
