package handlers

import (
	"dating-site-api/internal/models"
	"dating-site-api/pkg/services"
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

	if err := c.ShouldBindJSON(&userLogIn); err != nil {
		log.Fatal("LoginHandler: unable to bind data")
		return
	}

	// Достаем пользователя по имени из бд
	user, errQuery := services.GetUserByUsername(userLogIn.Username)
	if errQuery != nil {
		log.Printf("LoaginHandler: error during query proccesing: %s", errQuery.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "there is no such user"})
		return
	}

	// Сравниваем введённый парорль с паролем из бд
	errPassword := utils.ComparePbkdf2Sha256Hashes(user.Password.String, userLogIn.Password)
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
		return
	}

	// Добавляем токен в заголовок авторизации
	c.Header("Authorization", userToken)
	c.IndentedJSON(http.StatusOK, gin.H{"status": "Authorization completed"})
}

func (h Handler) RegisterHandler(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Printf("RegisterHandler: Error during bind json: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user format"})
		return
	}

	newUser.Password.String = utils.GeneratePbkdf2Sha256Hash(newUser.Password.String)

	// Зарос на создаие нового пользователя
	errQuery := services.CreateUser(newUser)
	if errQuery != nil {
		log.Printf("error during query processing: %s", errQuery.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There are some issues, try again later"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"status": "Registration is successfull"})
}
