package handlers

import (
	"dating-site-api/pkg/middlewares"
	"os"

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
			auth.POST("login/", h.LoaginHandler)
			auth.POST("register/", h.RegisterHandler)
		}

		user := api.Group("user")
		// Активием использование промежуточного программного компонента
		// для проверки jwt-token - ов, передавая секретный ключ для подписи
		user.Use(middlewares.JWTMiddleware(os.Getenv("JWT_SECRET_KEY")))
		{
			user.GET("list/", h.GetAllUsers)
			user.GET(":id", h.GetUserById)
			user.POST(":id/update", h.UpdateUserById)
		}
	}

	return router
}
