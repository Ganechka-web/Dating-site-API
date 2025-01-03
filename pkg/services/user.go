package services

import (
	"context"
	"dating-site-api/internal/database"
	"dating-site-api/internal/models"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

func GetAllActiveUsers() ([]models.User, error) {
	var users []models.User

	query := `SELECT id, age, username, email, city, 
			      date_birth, phone, description, gender 
			  FROM accounts_datinguser
			  WHERE is_active = true;`
	rows, errQuery := database.ConnectionPool.Query(context.Background(), query)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		errScan := rows.Scan(
			&user.ID, &user.Age, &user.Username, &user.Email, &user.City,
			&user.DateBirth, &user.Phone, &user.Description, &user.Gender,
		)
		if errScan != nil {
			log.Panicf("GetActiveUsers: error during query scan: %s", errScan.Error())
		}
		users = append(users, user)
	}

	return users, nil
}

func GetActiveUserById(userId int) (models.User, error) {
	var user models.User

	query := `SELECT id, age, username, email, city, 
			      date_birth, phone, description, gender 
			  FROM accounts_datinguser
			  WHERE is_active = true AND id = $1;`
	row := database.ConnectionPool.QueryRow(context.Background(), query, userId)

	errScan := row.Scan(
		&user.ID, &user.Age, &user.Username, &user.Email, &user.City,
		&user.DateBirth, &user.Phone, &user.Description, &user.Gender,
	)
	if errScan != nil {
		return user, errScan
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	query := `SELECT id, username, password
	          FROM accounts_datinguser
			  WHERE is_active = true AND username = $1`
	row := database.ConnectionPool.QueryRow(context.Background(), query, username)

	errScan := row.Scan(
		&user.ID, &user.Username, &user.Password,
	)
	if errScan != nil {
		return user, errScan
	}

	return user, nil
}

func CreateUser(newUser models.User) error {
	query := `INSERT INTO accounts_datinguser 
			      (age, username, email, password, city, date_birth, 
				   phone, description, gender, first_name, last_name, 
				   date_joined, is_superuser, is_staff, is_active)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, 
			  	      $13, $14, $15)`
	_, errQuery := database.ConnectionPool.Exec(context.Background(), query,
		newUser.Age, newUser.Username, newUser.Email, newUser.Password,
		newUser.City, newUser.DateBirth, newUser.Phone, newUser.Description,
		newUser.Gender, string(""), string(""), time.Now(), false, false, true,
	)
	if errQuery != nil {
		return errQuery
	}
	return nil
}

func UpdateUserById(userId int, data map[string]interface{}) error {
	params := make([]string, 0)

	// Валидируем и синхрогизируем типы для sql
	for key, value := range data {
		switch key {
		case "age":
			if age, ok := value.(float64); !ok {
				errText := fmt.Sprintln("can`t convert age: ", value)
				return errors.New(errText)

			} else {
				params = append(params, fmt.Sprintf("%s = %d", key, int8(age)))
			}
		case "date_birth":
			if dateBirth, ok := value.(string); !ok {
				errText := fmt.Sprintln("can`t convert date_birth: ", value)
				return errors.New(errText)

			} else {
				dateBirthParsed, errParse := time.Parse("02-01-2006", dateBirth)
				if errParse != nil {
					return errParse
				}
				params = append(params, fmt.Sprintf("%s = '%s'", key, dateBirthParsed.Format("02-01-2006")))
			}
		case "gender", "phone", "description", "city", "Email":
			params = append(params, fmt.Sprintf("%s = '%s'", key, value))
		default:
			errText := fmt.Sprintf("can`t convert field %s: %v", key, value)
			return errors.New(errText)
		}
	}
	query := fmt.Sprintf(`UPDATE accounts_datinguser
						  SET %s
						  WHERE id = $1;`, strings.Join(params, ", "))
	log.Printf("%s", query)
	_, errQuery := database.ConnectionPool.Exec(context.Background(), query, userId)
	if errQuery != nil {
		return errQuery
	}

	return nil
}
