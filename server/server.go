package server

import (
	"context"
	f "fmt"
	"net/http"
	"time"
)

type APIServer struct {
	// Создаём структуру сервера (просто надстройка над http.Server)
	httpServer *http.Server
}

func (as *APIServer) Run(port string, handler http.Handler) error {
	// Функция запускающая сервер, инкапсулируем все параметры запуска сервера.
	// Принимает порт на котором будет запущен сервер, а также хендлер.
	as.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	// бесконечный цикл для отслеживания запросов
	f.Println("Server`ve started and listening 8001 port")
	return as.httpServer.ListenAndServe()
}

func (as *APIServer) Shutdown(ctx context.Context) error {
	return as.httpServer.Shutdown(ctx)
}
