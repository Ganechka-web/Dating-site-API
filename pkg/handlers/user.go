package handlers

import (
	"dating-site-api/internal/models"
	"dating-site-api/pkg/services"
	"log"
	"net/http"
	"strconv"

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

	errQuery := services.UpdateUserById(userId, data)
	if errQuery != nil {
		log.Printf("UpdateUserById: error during query processing: %s", errQuery.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalied data format"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "user has been updated successfully"})
}
