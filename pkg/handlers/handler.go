package handlers

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	// создаем обработчик для маршрутов и запросов через Gin
	router := gin.New()

	api := router.Group("api/")
	{
		auth := api.Group("auth/")
		{
			auth.POST("register/", h.RegisterHandler)
			auth.POST("login/", h.LoaginHandler)
		}
	}

	return router
}
