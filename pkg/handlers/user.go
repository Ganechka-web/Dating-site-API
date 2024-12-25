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
	query := `SELECT id, age, username, email, city, 
			      date_birth, phone, description 
			  FROM accounts_datinguser;`
	rows, err := database.ConnectionPool.Query(context.Background(), query)

	if err != nil {
		log.Printf("Error during query processing %s", err.Error())
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Age, &user.Username, &user.Email, &user.City,
			&user.DateBirth, &user.Phone, &user.Description)
		if err != nil {
			log.Printf("Unable to read row %s", err.Error())
		}
		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (h Handler) GetUserById(c *gin.Context) {
	var user models.User
	userId := c.Param("id")

	query := `SELECT id, age, username, email, city, 
			      date_birth, phone, description 
			  FROM accounts_datinguser
			  WHERE id = $1;`
	row := database.ConnectionPool.QueryRow(context.Background(), query, userId)
	err := row.Scan(
		&user.ID, &user.Age, &user.Username, &user.Email, &user.City,
		&user.DateBirth, &user.Phone, &user.Description)

	if err != nil {
		log.Printf("Error during query scan: %s", err.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "There are no such user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
