package tests

import (
	"MEDODS/handlers"
	"MEDODS/pkg/repository"
	"MEDODS/pkg/service"
	"MEDODS/pkg/tokens"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Загрузка данных из config.yml
func initConfig() error {
	viper.AddConfigPath("../../config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// Что бы тесты работали нужно инициализировать все сервисы как в main
func initServices() *handlers.Handlers {
	//Загржаем данные из config.yml
	err := initConfig()
	if err != nil {
		logrus.Fatalf("initializing config error: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("dbPostgres.host"),
		Port:     viper.GetString("dbPostgres.port"),
		UserName: viper.GetString("dbPostgres.username"),
		DBName:   viper.GetString("dbPostgres.dbname"),
		SSLMode:  viper.GetString("dbPostgres.sslmode"),
		Password: viper.GetString("dbPostgres.password"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	token, _ := tokens.New(viper.GetString("signingKey"))
	services := service.NewService(repos, token)
	handlers := handlers.New(services)
	return handlers
}

// ТЕСТЫ НА /get_token (в test_table находятся сами тесты. Здесь они запускаются)
func TestHandler_getTokens(t *testing.T) {
	//Поднимаем все сервисы
	handlers := initServices()
	router := handlers.InitRoutes()

	//В цикле проходимся по каждому тесту и запускаем его
	for _, test := range getTestTableGetToken() {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/get_token", strings.NewReader(test.requestBody))

			w := httptest.NewRecorder()
			//Выполняется запрос
			router.ServeHTTP(w, req)

			//Проверяем ответ
			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}

// ТЕСТЫ НА /get_refresh (в test_table находятся сами тесты. Здесь они запускаются)
func TestHandler_refreh_token(t *testing.T) {
	//Поднимаем все сервисы
	handlers := initServices()
	router := handlers.InitRoutes()

	accessToken := ""
	refreshToken := ""

	//Сначала получаем токены
	test := getTestTableGetTokenStart()
	t.Run(test.name, func(t *testing.T) {
		//Для получение токенов
		type Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		var data Data

		req, _ := http.NewRequest("POST", "/get_token", strings.NewReader(test.requestBody))

		w := httptest.NewRecorder()
		//Выполняется запрос
		router.ServeHTTP(w, req)
		//Получаем ответ
		responseData, _ := io.ReadAll(w.Body)
		json.Unmarshal([]byte(responseData), &data)

		accessToken = data.AccessToken
		refreshToken = data.RefreshToken

		//Проверяем ответ
		assert.Equal(t, test.expectedCode, w.Code)
	})

	//Запускаем основные тесты
	//В цикле проходимся по каждому тесту и запускаем его
	for _, test := range getTestTableRefreshToken(accessToken, refreshToken) {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/refresh_token", strings.NewReader(test.requestBody))

			w := httptest.NewRecorder()
			//Выполняется запрос
			router.ServeHTTP(w, req)

			//Проверяем ответ
			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}
