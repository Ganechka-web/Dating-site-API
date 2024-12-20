package handlers

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	// создаем обработчик для маршрутов и запросов через Gin
	router := gin.New()

	// В ручную активируем логирование и обработку ошибок
	// с помощью gin middleware
	router.Use(gin.Logger(), gin.Recovery())

	api := router.Group("api/")
	{
		auth := api.Group("auth/")
		{
			auth.POST("register/", h.RegisterHandler)
			auth.POST("login/", h.LoaginHandler)
		}
		user := api.Group("user")
		{
			user.GET("list/", h.GetAllUsers)
		}
	}

	return router
}
