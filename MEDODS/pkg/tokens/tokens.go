package tokens

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager interface {
	ParseRefreshToken(refreshToken string) (string, string, string, bool, error)
	ParseAccessToken(accessToken string) (string, error)
	NewToken(ip, guid string, timeToLive int, status bool, refreshTokenUuid string) (string, error)
}

type Tokens struct {
	singingKey string
}

type tokenClaims struct {
	jwt.StandardClaims
	Ip               string
	Guid             string
	RefreshTokenUuid string
}

// Получаем секретный ключ
func New(singingKey string) (*Tokens, error) {
	//Проверяем не пустой ли ключ
	if singingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Tokens{singingKey: singingKey}, nil
}

// Создаём новый token
func (t *Tokens) NewToken(ip, guid string, timeToLive int, status bool, refreshTokenUuid string) (string, error) {
	expiredAt := time.Minute * time.Duration(timeToLive)
	//Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		//Создаём payload и указываем время жизни
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiredAt).Unix(),
		},
		ip,
		guid,
		refreshTokenUuid,
	})
	//Создаём токен
	tokenJWT, err := token.SignedString([]byte(t.singingKey))
	if err != nil {
		return "", err
	}

	//Статус true отвечает за то что это refresh token и его кодируем в base64
	if status {
		return base64.StdEncoding.EncodeToString([]byte(tokenJWT)), nil
	}

	//Возвращаем токен в base64
	return tokenJWT, nil
}

// Parse refresh token
func (t *Tokens) ParseRefreshToken(refreshToken string) (string, string, string, bool, error) {
	//Так как токен изначально в base64 нужно его декодировать
	decodedBytes, err := base64.URLEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", "", "", false, err
	}

	//Распарсиваем токен
	token, err := jwt.ParseWithClaims(string(decodedBytes), &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.singingKey), nil
	})

	//Если у токена закончился срок жизни. Сначало проверяем не err потому что если у
	//токена закончился срок жизни нужно вернуть то что пользователь не авторизован
	if !token.Valid {
		return "", "", "", true, err
	}
	//Если ошибка есть и она не звязана с тем что у токена закончился срок жизни
	if err != nil {
		return "", "", "", false, err
	}
	//Получаем значения из payload
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", "", "", false, errors.New("can`t get claims from refresh token")
	}

	return claims.Guid, claims.Ip, claims.RefreshTokenUuid, false, nil

}

// Parse access token, из access на надо получить refresh_uuid из payload чтобы удостовериться что токены связаны
func (t *Tokens) ParseAccessToken(accessToken string) (string, error) {
	//Распарсиваем токен
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.singingKey), nil
	})

	//Чтобы можно было получить данные из токена который истёк
	v, _ := err.(*jwt.ValidationError)

	//Если токен не валиден
	if !token.Valid {

		//Если ошибка в том что токен не валиден
		if v.Errors == jwt.ValidationErrorExpired {
			//Получаем значения из payload
			claims, ok := token.Claims.(*tokenClaims)
			if !ok {
				return "", errors.New("can`t get claims from access token")
			}
			return claims.RefreshTokenUuid, nil

		} else {
			return "", err
		}

	} else {
		//Получаем значения из payload
		claims, ok := token.Claims.(*tokenClaims)
		if !ok {
			return "", errors.New("can`t get claims from access token")
		}
		return claims.RefreshTokenUuid, nil
	}

}
