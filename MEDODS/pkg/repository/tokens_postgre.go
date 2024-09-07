package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// В этом файле функции для взаимодейтсвиями с tokens в БД
type TokensPostgre struct {
	db *sqlx.DB
}

func NewTokensPostgre(db *sqlx.DB) *TokensPostgre {
	return &TokensPostgre{db: db}
}

func (t *TokensPostgre) SaveRefreshToken(tokenRefresh, guid string) error {
	timeNow := time.Now().Format("2006-01-02 15:04:05")

	//Запрос на сохранение refresh token в БД в таблицу users
	query_ := fmt.Sprintf(`
	UPDATE %s 
	SET refresh_token = '%s', available = '%t', created_at = '%s' 
	WHERE user_id = (SELECT id FROM %s WHERE guid = '%s');
	`, refreshTokensTable, tokenRefresh, true, timeNow, usersTable, guid)

	//Исполняем запрос
	result, err := t.db.Exec(query_)
	if err != nil {
		return err
	}

	//Смотрим сколько строк обновили(Так как если неккоректные данные, он не кидает ошбку)
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("don`t update any rows.propably not correct data")
	}

	return nil
}

func (t *TokensPostgre) GetRefreshToken(guid string) (string, time.Time, error) {
	//Создай запрос на получение refresh token и expire_at
	query_ := fmt.Sprintf(`
		SELECT refresh_token, expire_at 
		FROM %s 
		WHERE refresh_tokens.user_id = (
			SELECT id FROM %s
			WHERE guid = '%s'
			);
		`, refreshTokensTable, usersTable, guid)

	//Переменная для получения token
	var token string
	//Полученик времени
	time := time.Now()

	//Исполняем запрос
	row := t.db.QueryRow(query_)
	if err := row.Scan(&token, &time); err != nil {
		return "", time, err
	}

	return token, time, nil
}

func (t *TokensPostgre) CheckAvailableToken(guid string) (bool, error) {
	//Создай запрос на получение available
	query_ := fmt.Sprintf(`
		SELECT available 
		FROM %s 
		WHERE refresh_tokens.user_id = (
			SELECT id FROM %s
			WHERE guid = '%s'
			);
		`, refreshTokensTable, usersTable, guid)

	//Переменная для получения статуса available
	var available bool

	//Исполняем запрос
	row := t.db.QueryRow(query_)
	if err := row.Scan(&available); err != nil {
		return available, err
	}

	return available, nil
}
