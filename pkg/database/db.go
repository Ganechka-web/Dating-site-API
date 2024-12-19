package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigDB struct {
	DBName         string
	DBUser         string
	DBUserPassword string
	DBHost         string
	DBPort         string
}

// Синглтон для подключение к бд (игнорируем внедрение зависимотей)
var ConnectionPool *pgxpool.Pool

func ConnectToDB(ctx context.Context, config *ConfigDB) *pgxpool.Pool {
	dbconn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.DBHost, config.DBUser, config.DBUserPassword,
		config.DBName, config.DBPort,
	)

	var err error
	ConnectionPool, err = pgxpool.New(ctx, dbconn)
	if err != nil {
		// В случае ошибки пишем в канал вывода ошибок
		// Аварийно завершаем работу
		fmt.Fprintf(os.Stderr, "Uneble to connect to database %v", err)
		os.Exit(1)
	}

	return ConnectionPool
}
