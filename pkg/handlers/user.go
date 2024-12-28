package handlers

import (
	"context"
	"dating-site-api/internal/database"
	"dating-site-api/internal/models"
	"dating-site-api/pkg/services"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAllUsers(c *gin.Context) {
	users, err := services.GetAllActiveUsers()
	if err != nil {
		log.Printf("GetAllUsers: error during query processing: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There are some issues, try again later"})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (h Handler) GetUserById(c *gin.Context) {
	var user models.User

	userId, errAtoi := strconv.Atoi(c.Param("id"))
	if errAtoi != nil {
		log.Printf("UpdateUserById: error during param convert: %s", errAtoi.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	user, errQuery := services.GetActiveUserById(userId)
	if errQuery != nil {
		log.Printf("Error during query proccesing: %s", errQuery.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "There are no such user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (h Handler) UpdateUserById(c *gin.Context) {
	// Карта для приёма переменного количества аргументов которых нужно обновить
	var data map[string]interface{}
	var params []string

	userId, errAtoi := strconv.Atoi(c.Param("id"))
	if errAtoi != nil {
		log.Printf("UpdateUserById: error during param convert: %s", errAtoi.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	if errBind := c.ShouldBindJSON(&data); errBind != nil {
		log.Printf("UpdateUserById: error during binding: %s", errBind.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
		return
	}

	// Формируем формат для UPDATE запроса, учитывая типы данных
	for key, value := range data {
		switch key {
		case "age":
			// Проверяем и изменяем стандартный тип float64
			if age, ok := value.(float64); !ok {
				log.Println("UpdateUserById: Invalid age format")
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
				return
			} else {
				params = append(params, fmt.Sprintf("%s = %d", key, uint8(age)))
			}
		case "date_birth":
			// Проверяем корректность формата полученной даты
			if dateBirth, ok := value.(string); !ok {
				log.Println("UpdateUserById: Invalid date format")
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
				return
			} else {
				dateBirthParsed, errDate := time.Parse("2006-01-02", dateBirth)
				if errDate != nil {
					log.Println("UpdateUserById: Invalid date format")
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
					return
				}
				// Преобразуем дату в формат совместмый с SQL
				params = append(params, fmt.Sprintf("%s = '%s'", key, dateBirthParsed.Format("2006-01-02")))
			}
		default:
			params = append(params, fmt.Sprintf("%s = '%s'", key, value))
		}
	}

	query := fmt.Sprintf(`UPDATE accounts_datinguser
						  SET %s
						  WHERE id = $1;`, strings.Join(params, ", "))
	fmt.Println(query)
	_, errQuery := database.ConnectionPool.Exec(context.Background(), query, userId)
	if errQuery != nil {
		log.Printf("UpdateUserById: error during query processing: %s", errQuery.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There ae some issues, try again later"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "user has been updated successfully"})
}
