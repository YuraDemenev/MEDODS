package tokens

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	NewJWT(userGUID string) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type Tokens struct {
	singingKey string
}

// Time to live for JWT token
const timeToLive = time.Minute * 15

func New(singingKey string) (*Tokens, error) {
	if singingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Tokens{singingKey: singingKey}, nil
}

// Create JWT token
func (t *Tokens) NewJWT(userGUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(timeToLive).Unix(),
		Subject:   userGUID,
	})
	return token.SignedString([]byte(t.singingKey))
}

// Parse JWT token
func (t *Tokens) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.singingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

// Create a random refresh token
func (t *Tokens) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	str := fmt.Sprintf("%x", b)
	tokenBase64 := base64.StdEncoding.EncodeToString([]byte(str))

	return tokenBase64, nil
}
