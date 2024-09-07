package main

import (
	"MEDODS/elements"
	"MEDODS/handlers"
	"MEDODS/pkg/repository"
	"MEDODS/pkg/service"
	"MEDODS/pkg/tokens"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	//Загружаем конфиг(config.yml)
	err := initConfig()
	//Проверяем наличие ошибки при загрузки
	if err != nil {
		logrus.Fatalf("initializing config error: %s", err.Error())
	}
	fmt.Println("Succes config load")

	//Подключаемся к PostgreSQL
	db, err := repository.NewPostgresDB(repository.Config{
		//Подстваляем данные из config
		Host:     viper.GetString("dbPostgres.host"),
		Port:     viper.GetString("dbPostgres.port"),
		UserName: viper.GetString("dbPostgres.username"),
		DBName:   viper.GetString("dbPostgres.dbname"),
		SSLMode:  viper.GetString("dbPostgres.sslmode"),
		Password: viper.GetString("dbPostgres.password"),
	})

	//Проверяем наличие ошибки
	if err != nil {
		logrus.Fatalf("Cant connect to postgreSQL. Err: %s", err.Error())
	}
	fmt.Println("Succes connect to postgreSQL")

	//Инициализируем репозиторий(Здесь находсятся все функции, которые будут взаимодействовать с БД).
	//Функции будут запаскаться через services
	repos := repository.NewRepository(db)
	//В token будут создоваться и обрабатываться токены
	token, err := tokens.New(viper.GetString("signingKey"))
	if err != nil {
		logrus.Fatalf("Problem with signing key. Err:%s", err)
	}

	//Инициализируем сервисы(Сервисы будут использоваться для запуска, если понадобиться дополнительных функций
	//перед взаимодействием с БД и через сервисы будут запускаться функции для взаимодействия с БД).
	//Функции сервисов запускаются из Handlers
	services := service.NewService(repos, token)

	//Инициализируем handler(здесь мы задаём пути запроса и также указываем какие функции запустяться по данному пути)
	handlers := handlers.New(services)
	//Создаем сервер
	server := new(elements.Server)
	//Запускаем сервер
	err = server.Run(viper.GetString("port"), handlers.InitRoutes())
	if err != nil {
		logrus.Fatalf("error while runnig server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
