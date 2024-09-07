package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// В этом файле функции для взаимодейтсвиями с tokens в БД
type UserPostgre struct {
	db *sqlx.DB
}

func NewUserPostgre(db *sqlx.DB) *UserPostgre {
	return &UserPostgre{db: db}
}

// Функция на получение почты пользователя через guid
func (u *UserPostgre) GetUserEmail(guid string) (string, error) {
	//Запрос на получение почты user для отправки warning
	query_ := fmt.Sprintf(`
	SELECT email FROM %s WHERE guid=$1
	`, usersTable)

	//Переменная для получения почты
	var email string

	//Исполняем запрос
	row := u.db.QueryRow(query_, guid)
	if err := row.Scan(&email); err != nil {
		return "", err
	}
	return email, nil
}
