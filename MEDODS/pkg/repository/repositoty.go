package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Tokens
	User
}

type Tokens interface {
	SaveRefreshToken(refreshToken string, guid string) error
	GetRefreshToken(guid string) (string, time.Time, error)
	CheckAvailableToken(guid string) (bool, error)
}

type User interface {
	GetUserEmail(guid string) (string, error)
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Tokens: NewTokensPostgre(db),
		User:   NewUserPostgre(db),
	}
}
