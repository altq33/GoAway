package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *Handler) *gin.Engine {
	router := gin.Default()

	h.InitRoutes(router)

	return router
}
