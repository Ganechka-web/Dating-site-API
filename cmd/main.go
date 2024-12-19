package main

import (
	"context"
	"dating-site-api/pkg/database"
	"dating-site-api/pkg/handlers"
	"dating-site-api/server"

	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Подгружаем переменные окружения
	err_dotenv := godotenv.Load(".env")
	if err_dotenv != nil {
		fmt.Fprintf(os.Stderr, "Eroor loading .env file: %v", err_dotenv)
		os.Exit(1)
	}

	// Точка входа в приложение, создвем новый хендлер для запросов,
	// запускаем сервер на 8001 порту с новым хендлером
	handlers := new(handlers.Handler)

	// Подключаемся к бд
	config := database.ConfigDB{
		DBName: os.Getenv("DB_NAME"), DBUser: os.Getenv("DB_USER"),
		DBUserPassword: os.Getenv("DB_USER_PASSWORD"),
		DBHost:         os.Getenv("DB_HOST"), DBPort: os.Getenv("DB_PORT")}
	connection := database.ConnectToDB(context.Background(), &config)

	defer connection.Close()

	// Создаём и запускаем сервер
	Server := server.APIServer{}
	if err := Server.Run("8001", handlers.InitRoutes()); err != nil {
		log.Fatalf("There are some issues during server starts")
	}
}
