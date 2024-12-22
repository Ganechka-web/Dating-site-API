package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Решено сделать авторизизацию через jwt-token
// Создаем middleware для проверки подлиности и подписи jwt-token
func JWTMiddleware(secret_key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			// В случае отсуцтвия токена в headers пррываем работы функции обёртки и основной
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "You aren`t logged in"})
			c.Abort()
			return
		}

		// Парсим токен в байтовый срез
		tokenByte, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Проверям подлинность метода подписи в случае успеха возвращаем ключ подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid claim method: %v", token.Header["alg"])
			}
			// Возвращаем срез байтов на основе секретного ключа
			return []byte(secret_key), nil
		})

		if err != nil || !tokenByte.Valid {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			log.Printf("Invalid token %s", err.Error())
		}

		// Добавляем информацию о токене из claims в сонтекст gin
		if claims, ok := tokenByte.Claims.(jwt.MapClaims); ok {
			c.Set("id", claims["id"])
			c.Set("iat", claims["iat"])
			c.Set("exp", claims["exp"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claim"})
			c.Abort()
			return
		}

		c.Next()
	}
}
