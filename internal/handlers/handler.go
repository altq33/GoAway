package handlers

import (
	"goawayauth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.SecretService
}

func NewHandler(s service.SecretService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1/secrets")
	{
		api.POST("/", h.CreateSecret)
		api.GET("/:id/meta", h.GetSecretMeta)
		api.POST("/:id/read", h.ReadSecret)
	}
}

func (h *Handler) CreateSecret(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Ручка создания записки еще не готова",
	})
}

func (h *Handler) GetSecretMeta(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Ручка метаданных еще не готова, запрошен ID: " + id,
	})
}

func (h *Handler) ReadSecret(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Ручка чтения еще не готова, запрошен ID: " + id,
	})
}
