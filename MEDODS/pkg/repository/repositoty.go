package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Authorization
	RefreshTokens
}

type Authorization interface {
	CreateUser(userName, password string) (string, error)
	GetUserId(username, password string) (string, error)
	SaveRefreshToken(refreshToken string, guid string) (string, error)
	GetUserName(guid string) (string, error)
}

type RefreshTokens interface {
	CheckExpireTime(tokenRefresh, tokenRefreshGUID string) (string, error)
}

func New(client *mongo.Client) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(client),
		RefreshTokens: NewRefreshMongo(client),
	}
}
