package service

import (
	"MEDODS/pkg/repository"
	"MEDODS/pkg/tokens"
	"crypto/sha1"
	"errors"
	"fmt"
)

// For secure
const (
	salt = "123sad)__31ld{as}dk21"
)

type AuthService struct {
	repo  repository.Authorization
	token tokens.TokenManager
}

func NewAuthService(repo repository.Authorization, token tokens.TokenManager) *AuthService {
	return &AuthService{repo: repo, token: token}
}

// Generate Token
func (s *AuthService) GenerateToken(username, password string) (string, string, string, error) {
	//Get Guid from db
	guid, err := s.repo.GetUserId(username, generatePasswordHash(password))
	if err != nil {
		return "", "", "", err
	}
	if guid == "" {
		return "", "", "", errors.New("no such user")
	}

	//Create tokens
	jwtToken, err := s.token.NewJWT(guid)
	if err != nil {
		return "", "", "", err
	}

	refresthToken, err := s.token.NewRefreshToken()
	if err != nil {
		return "", "", "", err
	}

	//Save refreshToken to db
	refresthTokenGUID, err := s.repo.SaveRefreshToken(refresthToken, guid)
	if err != nil {
		return "", "", "", err
	}

	return jwtToken, refresthToken, refresthTokenGUID, nil
}

// Parse Token
func (s *AuthService) ParseToken(token string) (string, error) {
	guid, err := s.token.Parse(token)
	if err != nil {
		return "", err
	}

	return guid, nil
}

// Create User
func (s *AuthService) CreateUser(username, password string) (string, error) {
	return s.repo.CreateUser(username, generatePasswordHash(password))
}

// Get User name
func (s *AuthService) GetUserName(guid string) (string, error) {
	return s.repo.GetUserName(guid)
}

// Hash password
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
