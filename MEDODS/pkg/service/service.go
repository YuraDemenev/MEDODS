package service

import (
	"MEDODS/pkg/repository"
	"MEDODS/pkg/tokens"
)

type Service struct {
	Authorization
	RefreshTokens
}

type Authorization interface {
	GenerateToken(username, password string) (string, string, string, error)
	ParseToken(token string) (string, error)
	CreateUser(username, password string) (string, error)
	GetUserName(guid string) (string, error)
}

type RefreshTokens interface {
	RefreshTokens(tokenRefresh string, tokenRefreshGUID string) (string, string, string, string, error)
}

func NewService(repos *repository.Repository, token *tokens.Tokens) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, token),
		RefreshTokens: NewRefreshTokensService(repos.RefreshTokens, repos.Authorization, token),
	}
}
