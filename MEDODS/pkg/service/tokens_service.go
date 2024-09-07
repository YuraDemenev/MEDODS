package service

import (
	"MEDODS/pkg/repository"
	"MEDODS/pkg/tokens"
	"errors"

	"github.com/google/uuid"
)

// Устанавливаем время жизни access токена
const timeToLiveAcessToken = 3

// Устанавливаем время жизни refresh токена
const timeToLiveRefreshToken = 5

type TokensService struct {
	repoToken repository.Tokens
	repoUser  repository.User
	token     tokens.TokenManager
}

func NewTokensService(repoToken repository.Tokens, token tokens.TokenManager, repoUser repository.User) *TokensService {
	return &TokensService{repoToken: repoToken, token: token, repoUser: repoUser}
}

// Generate Token
func (t *TokensService) GenerateNewTokens(userGuid, ip string) (string, string, error) {
	//Возвращаем наши токены
	return createRefreshAndAccessTokensAndSaveInBD(userGuid, ip, t)
}

// Обновляем токены
func (t *TokensService) RefreshTokens(accessToken, refreshToken string, ipLast string) (string, string,
	string, bool, error) {
	// Парсим refresh token, получаем ip и userGuid из payload, и проверяем не истёк ли срок действия
	userGuid, ip, refreshUUID, checkExpire, err := t.token.ParseRefreshToken(refreshToken)
	//Если вернулся true значить у токена закончился срок жизни
	if checkExpire {
		return "", "", "", checkExpire, err
	}

	//Парсим access token чтобы удостовериться что токены были выданы вместе
	refreshUUIDFromAccessToken, err := t.token.ParseAccessToken(accessToken)
	if err != nil {
		return "", "", "", false, err
	}

	//Если у токенов не равны UUID значит они выданы не вместе
	if refreshUUID != refreshUUIDFromAccessToken {
		return "", "", "", false, errors.New("tokens have diffrent uuids")
	}

	email := ""
	//Сравниваем ip предыдущий с нынешним
	if ip != ipLast {
		//Если не равно нужно послать на почту warning, для этого
		//получаем из БД почту пользователя через userGuid
		email, err = t.repoUser.GetUserEmail(userGuid)
		if err != nil {
			return "", "", "", false, err
		}
	}

	//Проверяем можно ли использовать этот refresh token (В БД в таблице refresh_tokens есть поле available,
	//оно отвечает за то можно ли использовать токен или нет и проверяем это поле)
	checkAvailable, err := t.repoToken.CheckAvailableToken(userGuid)
	if err != nil {
		return "", "", email, false, err
	}
	//Проверяем не забанен ли этот токен
	if !checkAvailable {
		return "", "", email, false, errors.New("this token is not available")
	}

	//Делаем перевыпуск токенов
	refreshToken, accessToken, err = createRefreshAndAccessTokensAndSaveInBD(userGuid, ip, t)
	if err != nil {
		return "", "", email, false, err
	}

	//Возвращаем токены
	return refreshToken, accessToken, email, false, nil

}

// Функция для создание refresh, access tokens и сохранение refresh в БД.
// Так как мы дважды создаём токены, вынесем это в отдельную функцию
func createRefreshAndAccessTokensAndSaveInBD(userGuid, ip string, t *TokensService) (string, string, error) {
	//Чтобы токены были связаны у мы создаём uuid у access и refresh оно одинаковое
	refreshTokenUuid := uuid.New().String()

	//Создаём новый access token (false показывает что это access token)
	accessToken, err := t.token.NewToken(ip, userGuid, timeToLiveAcessToken, false, refreshTokenUuid)
	if err != nil {
		return "", "", err
	}

	//Создаём к нему refresh token (true показывает что это refresh token)
	refreshToken, err := t.token.NewToken(ip, userGuid, timeToLiveRefreshToken, true, refreshTokenUuid)
	if err != nil {
		return "", "", err
	}

	//Сохраняем токен в БД
	err = t.repoToken.SaveRefreshToken(refreshToken, userGuid)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
