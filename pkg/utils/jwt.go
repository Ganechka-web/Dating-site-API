package utils

import (
	"dating-site-api/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Подгружаем переменные окружения
func LoadDotEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
}

// Генерируем новый jwt-token по id и текущему времени,
// вощвращаем подписывая токен секретным ключём
func GenerateJWTToken(user models.User) (string, error) {
	LoadDotEnv()

	// Генирируем токен по id пользователя сроком 1 день
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
