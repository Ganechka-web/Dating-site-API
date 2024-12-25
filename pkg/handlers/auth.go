package handlers

import (
	"context"
	"dating-site-api/internal/database"
	"dating-site-api/internal/models"
	"dating-site-api/pkg/utils"
	"log"
	"net/http"
	"time"

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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Your username or password didn`t match"})
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

	// (is_superuser, first_name, last_name, date_joined, is_staff, is_active) для совместимости с django ORM и бд
	query := `INSERT INTO accounts_datinguser (age, username, email, password,
			      city, date_birth, phone, description, gender, first_name, last_name, 
				  date_joined, is_superuser, is_staff, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, errQuery := database.ConnectionPool.Exec(context.Background(), query,
		newUser.Age, newUser.Username, newUser.Email, utils.GeneratePbkdf2Sha256Hash(newUser.Password.String),
		newUser.City, newUser.DateBirth, newUser.Phone, newUser.Description, newUser.Gender, string(""), string(""), time.Now(), false, false, true)
	if errQuery != nil {
		log.Printf("error during query processing: %s", errQuery.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There are some issues, try again later"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"status": "Registration is successfull"})
}
