package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Названия таблиц в базе данных
const usersTable = "users"
const refreshTokensTable = "refresh_tokens"

// Структура Config с информацией для подключения к БД
type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.UserName, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
