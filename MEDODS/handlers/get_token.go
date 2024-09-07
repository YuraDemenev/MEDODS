package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Структура для получения json из post запроса
type getTokenStruct struct {
	Guid string `json:"guid"`
	Ip   string `json:"ip"`
}

// Структура для отправки токенов в ввиде json
type setTokenStruct struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

type SetErrorStuct struct {
	Error string `json:"error"`
	Info  string `json:"info"`
	//Для того чтобы знать на какую почту прислать warning, что ip поменялось
	Email string `json:"email"`
}

// Словарь чтобы просто и быстро получать строку ответа по коду ошибки
var DictionaryErrors = map[int]string{
	400: "{'error':' Неккорентные данные'}",
	401: "{'error':' Пользователь не авторизован'}",
	403: "{'error':' Пользователь не имеет доступа'}",
	500: "{'error':' Внутренняя ошибка сервера'}",
}

// Функция на получение Access и Refresh Token
func (h *Handlers) getToken(c *gin.Context) {
	//Структура для получения данных
	var requestData getTokenStruct

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

	//Проверяем не пустой ли  guid/ip
	if requestData.Guid == "" || requestData.Ip == "" {
		//Возвращаем ошибку
		ReturnError(c, "Guid or ip is empty", http.StatusBadRequest, "")
		return
	}

	//Получаем токены
	accessToken, refreshToken, err := h.service.Tokens.GenerateNewTokens(requestData.Guid, requestData.Ip)
	if err != nil {
		str := fmt.Sprintf("Can`t generate token. Error: %s", err.Error())
		//Возвращаем ошибку
		ReturnError(c, str, http.StatusInternalServerError, "")
		return
	}

	//Создаём структуру ответа
	var setTokenStruct setTokenStruct
	setTokenStruct.Access_token = accessToken
	setTokenStruct.Refresh_token = refreshToken

	//Конвертруем в json  заранее заготовленную структуры, для возвращения ответа
	result, err := json.Marshal(setTokenStruct)
	if err != nil {
		str := fmt.Sprintf("Can`t convert setTokenStuct to json. Error: %s", err.Error())
		//Возвращаем ошибку
		ReturnError(c, str, http.StatusInternalServerError, "")
		return
	}

	//Возвращаем токены пользователю
	c.Writer.Write(result)
}

// Функция для возвращения ошибки
func ReturnError(c *gin.Context, str string, errorNum int, email string) {
	logrus.Error(str)
	//Возвращаем ошибку в postman

	//Создаём структуру для возвращения ошибки
	var setErrorStuct SetErrorStuct
	setErrorStuct.Error = DictionaryErrors[errorNum]
	setErrorStuct.Info = str

	//Добавляем email если он не пустой
	if email != "" {
		setErrorStuct.Email = email
	}

	//Конвертируем в json
	result, err := json.Marshal(setErrorStuct)
	if err != nil {
		str := fmt.Sprintf("Can`t convert setErrorStuct to json. Error: %s", err.Error())
		logrus.Error(str)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Возвращаем ошибку
	c.Writer.WriteHeader(errorNum)
	c.Writer.Write(result)
}
