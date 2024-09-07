package tests

import (
	"encoding/json"
	"net/http"
)

type TestTableStruct = struct {
	name         string
	requestBody  string
	expectedCode int
}

// Функция возвращает 1 тест для /get_token: Tест на получение токенов
// (Надо чтобы получить первые токены, для следующих тестов)
func getTestTableGetTokenStart() TestTableStruct {
	var testTableGetToken TestTableStruct
	testTableGetToken.name = "Test #1. OK. корректный стандартный запрос"
	testTableGetToken.requestBody = `{
	 			"guid":"b758e00b-9c42-4f2b-a84d-0af47b937d17",
     			"ip":"1.2.3"
	 		}`
	testTableGetToken.expectedCode = http.StatusOK

	return testTableGetToken
}

// Тесты на /get_token
func getTestTableGetToken() []TestTableStruct {
	testTableGetToken := []TestTableStruct{
		{
			name: "Test #1. OK. Корректный стандартный запрос",
			requestBody: `{
				"guid":"b758e00b-9c42-4f2b-a84d-0af47b937d17",
    			"ip":"1.2.3"
			}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test #2. NOT OK. Неверный guid",
			requestBody: `{
				"guid":"-1",
				"ip":"1.2.3"
			}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Test #3. NOT OK. Неверное request body",
			requestBody: `{
				"ababab":"babab"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Test #4. NOT OK. пустое guid",
			requestBody: `{
				"guid":"",
				"ip":"1.2.3"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Test #5. NOT OK. пустое ip",
			requestBody: `{
				"guid":"asd",
				"ip":""
			}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	return testTableGetToken
}

// Тесты на /refresh_token
func getTestTableRefreshToken(accessToken, refreshToken string) []TestTableStruct {
	type RequestBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Ip           string `json:"ip"`
	}
	var requestBody RequestBody
	requestBody.AccessToken = accessToken
	requestBody.RefreshToken = refreshToken
	requestBody.Ip = "1.2.3"

	//Конвертруем в json  заранее заготовленную структуры, для возвращения ответа
	requestBodyResult2, _ := json.Marshal(requestBody)
	//Для работы тестов меняем значения и создаем новые переменные
	//ТЕСТ #3
	requestBody.AccessToken = ""
	requestBodyResult3, _ := json.Marshal(requestBody)
	//ТЕСТ #4
	requestBody.AccessToken = accessToken
	requestBody.Ip = "1.1.1"
	requestBodyResult4, _ := json.Marshal(requestBody)
	//ТЕСТ #5
	requestBody.AccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	requestBodyResult5, _ := json.Marshal(requestBody)
	//ТЕСТ #6
	requestBody.AccessToken = accessToken
	requestBody.RefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	testTableGetToken := []TestTableStruct{
		{
			name:         "Test #2. OK. Корректный стандартный запрос",
			requestBody:  string(requestBodyResult2),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Test #3. NOT OK. Access token пустой",
			requestBody:  string(requestBodyResult3),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Test #4. OK. ip другой ",
			requestBody:  string(requestBodyResult4),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Test #5. NOT OK. access token неккорктный ",
			requestBody:  string(requestBodyResult5),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Test #6. NOT OK. refresh token неккорктный ",
			requestBody:  string(requestBodyResult5),
			expectedCode: http.StatusInternalServerError,
		},
	}

	return testTableGetToken
}
