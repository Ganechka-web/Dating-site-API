package main

import (
	"dating-site-api/pkg/handlers"
	"dating-site-api/server"
	"log"
)

func main() {
	// Точка входа в приложение, создвем новый хендлер для запросов,
	// запускаем сервер на 8001 порту с новым хендлером
	handlers := new(handlers.Handler)

	Server := server.APIServer{}
	if err := Server.Run("8001", handlers.InitRoutes()); err != nil {
		log.Fatalf("There are some issues during server starts")
	}
}
