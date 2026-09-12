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
	var req CreateSecretRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Некорректный запрос",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	serviceReq := service.CreateSecretDTO{
		Text:              req.Text,
		IsClientEncrypted: req.IsClientEncrypted,
		Password:          req.Password,
		TTLHours:          req.TTLHours,
	}

	result, err := h.service.CreateSecret(ctx, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось создать записку, попробуйте позже",
		})
		return
	}

	resp := CreateSecretResponse{
		ID:            result.ID,
		EncryptionKey: result.EncryptionKey,
	}

	c.JSON(http.StatusCreated, resp)
}

// GetSecretMeta возвращает информацию о том, нужен ли пароль и ключ
func (h *Handler) GetSecretMeta(c *gin.Context) {
	id := c.Param("id")

	meta, err := h.service.GetSecretMeta(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SecretMetaResponse{
		IsClientEncrypted: meta.IsClientEncrypted,
		HasPassword:       meta.RequiresPassword,
	})
}

// ReadSecret проверяет пароль, расшифровывает и сжигает записку
func (h *Handler) ReadSecret(c *gin.Context) {
	id := c.Param("id")

	var req ReadSecretRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Некорректный запрос",
			"details": err.Error(),
		})
		return
	}

	dto := service.ReadSecretDTO{
		ID:            id,
		Password:      req.Password,
		EncryptionKey: req.EncryptionKey,
	}

	text, err := h.service.ReadSecret(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ReadSecretResponse{
		Text: text,
	})
}
