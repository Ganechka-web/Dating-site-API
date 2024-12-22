package handlers

import (
	"context"
	"dating-site-api/internal/database"
	"dating-site-api/internal/models"
	"dating-site-api/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLogIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h Handler) LoaginHandler(c *gin.Context) {
	var userLogIn UserLogIn
	var user models.User

	if err := c.ShouldBindJSON(&userLogIn); err != nil {
		log.Fatal("LoginHandler: unable to bind data")
		return
	}

	// Достаем пользователя по имени из бд
	query := `SELECT id, username, password 
			  FROM accounts_datinguser 
			  WHERE username = $1;`
	row := database.ConnectionPool.QueryRow(context.Background(), query, userLogIn.Username)

	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		log.Printf("LogInHandler: unable to scan query: %v", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There some issues, try again later"})
	}

	// Сравниваем введённый парорль с паролем из бд
	errPassword := utils.ComparePbkdf2Sha265Hashes(user.Password.String, userLogIn.Password)
	if !errPassword {
		log.Println("LoginHandler: Password didn`t match")
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Password didn`t match"})
		return
	}

	// Генерируем токер авторизации для пользователя
	userToken, errToken := utils.GenerateJWTToken(user)
	if errToken != nil {
		log.Printf("LogInHandler: error with generating token %v", errToken.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There some issues, try again later"})
	}

	// Добавляем токен в заголовок авторизации
	c.Header("Authorization", userToken)
	c.IndentedJSON(http.StatusOK, gin.H{"status": "Authorization completed"})
}
