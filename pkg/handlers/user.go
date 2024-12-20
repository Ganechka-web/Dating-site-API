package handlers

import (
	"context"
	"dating-site-api/internal/database"
	"dating-site-api/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) GetAllUsers(c *gin.Context) {
	query := "SELECT id, age, username, email, city, date_birth, phone, description FROM accounts_datinguser;"
	rows, err := database.ConnectionPool.Query(context.Background(), query)

	if err != nil {
		log.Fatalf("Error during query processing %s", err.Error())
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Age, &user.Username, &user.Email, &user.City,
			&user.DateBirth, &user.Phone, &user.Description)
		if err != nil {
			log.Fatalf("Unable to read row %s", err.Error())
		}
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}
