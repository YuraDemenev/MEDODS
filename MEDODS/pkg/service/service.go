package service

import (
	"MEDODS/pkg/repository"
	"MEDODS/pkg/tokens"
)

type Service struct {
	Tokens
}

type Tokens interface {
	GenerateNewTokens(guid, ip string) (string, string, error)
	RefreshTokens(accessToken, refreshToken string, ip string) (string, string, string, bool, error)
}

func NewService(repos *repository.Repository, token *tokens.Tokens) *Service {
	return &Service{
		Tokens: NewTokensService(repos.Tokens, token, repos.User),
	}
}
