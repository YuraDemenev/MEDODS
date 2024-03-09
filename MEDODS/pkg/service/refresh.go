package service

import (
	"MEDODS/pkg/repository"
	"MEDODS/pkg/tokens"
)

type RefreshTokensService struct {
	repo     repository.RefreshTokens
	repoAuth repository.Authorization
	token    tokens.TokenManager
}

func NewRefreshTokensService(repo repository.RefreshTokens, repoAuth repository.Authorization, token tokens.TokenManager) *RefreshTokensService {
	return &RefreshTokensService{repo: repo, repoAuth: repoAuth, token: token}
}

func (r *RefreshTokensService) RefreshTokens(tokenRefresh string, tokenRefreshGUID string) (string, string, string, string, error) {
	guid, err := r.repo.CheckExpireTime(tokenRefresh, tokenRefreshGUID)
	if err != nil {
		return "", "", "", "", err
	}

	//Create tokens
	jwtToken, err := r.token.NewJWT(guid)
	if err != nil {
		return "", "", "", "", err
	}

	refresthToken, err := r.token.NewRefreshToken()
	if err != nil {
		return "", "", "", "", err
	}

	//Save refreshToken to db
	tokenRefreshGUID, err = r.repoAuth.SaveRefreshToken(refresthToken, guid)
	if err != nil {
		return "", "", "", "", err
	}

	return jwtToken, refresthToken, tokenRefreshGUID, guid, nil
}
