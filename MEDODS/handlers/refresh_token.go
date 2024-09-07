package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Структура для получения json из post запроса
type getRefreshTokenStruct struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Ip           string `json:"ip"`
}

// Структура для положительного ответа на refresh_tokens
type sendTokensStruct struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
}

// Функция на получение Access и Refresh Token
func (h *Handlers) refreshToken(c *gin.Context) {
	//Структура для получения данных
	var requestData getRefreshTokenStruct

	//Создаем декодер
	decoder := json.NewDecoder(c.Request.Body)
	//Запрещаем иметь в запросе поля которых нет в структуре
	decoder.DisallowUnknownFields()

	//Декодируем
	err := decoder.Decode(&requestData)
	if err != nil {
		str := fmt.Sprintf("Ошибка при декодирования request body. Ошибка: %s", err.Error())
		//Возвращаем ошибку
		ReturnError(c, str, http.StatusBadRequest, "")
		return
	}

	//Проверяем не пустой ли  refresh_token
	if requestData.RefreshToken == "" || requestData.AccessToken == "" {
		//Возвращаем ошибку
		ReturnError(c, "access/refresh token is empty", http.StatusBadRequest, "")
		return
	}

	//Обновляем токены
	accessToken, refreshToken, email, checkLiveToken, err := h.service.Tokens.RefreshTokens(requestData.AccessToken,
		requestData.RefreshToken, requestData.Ip)
	if err != nil {
		str := fmt.Sprintf("Can`t generate token. Error: %s", err.Error())
		//Возвращаем ошибку
		ReturnError(c, str, http.StatusInternalServerError, email)
		return
	}
	//Если срок жизни refresh token истёк
	if checkLiveToken {
		//Возвращаем ошибку что пользователь не авторизован
		str := "Refresh token expired"
		ReturnError(c, str, http.StatusUnauthorized, email)
		return
	}

	//Возвращаем ответ
	var sendTokensStruct sendTokensStruct
	sendTokensStruct.AccessToken = accessToken
	sendTokensStruct.RefreshToken = refreshToken
	sendTokensStruct.Email = email

	//Конвертруем в json  заранее заготовленную структуры, для возвращения ответа
	result, err := json.Marshal(sendTokensStruct)
	if err != nil {
		str := fmt.Sprintf("Can`t convert setTokenStuct to json. Error: %s", err.Error())
		//Возвращаем ошибку
		ReturnError(c, str, http.StatusInternalServerError, "")
		return
	}

	//Возвращаем токены пользователю
	c.Writer.Write(result)
}
