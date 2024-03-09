package main

import (
	"MEDODS/elements"
	"MEDODS/handlers"
	"MEDODS/pkg/repository"
	"MEDODS/pkg/service"
	"MEDODS/pkg/tokens"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	//load config
	err := initConfig()
	if err != nil {
		logrus.Fatalf("initializing config error: %s", err.Error())
	}
	fmt.Println("Succes config load")

	ctx := context.Background()

	//Connect to Mongo
	client, err := repository.NewMongoDB(repository.ConfigDB{
		Url:    viper.GetString("mongo.Url"),
		Addres: viper.GetString("mongo.Addres"),
	}, ctx)

	defer client.Disconnect(ctx)

	if err != nil {
		logrus.Fatalf("Cant connect to mongo. Err: %s", err.Error())
	}
	fmt.Println("Succes connect to mongo")

	//Initilize services
	repos := repository.New(client)
	token, err := tokens.New(viper.GetString("signingKey"))
	if err != nil {
		logrus.Fatalf("Problem with signing key. Err:%s", err)
	}

	services := service.NewService(repos, token)

	//Initilize handlers
	handlers := handlers.New(services)
	server := new(elements.Server)
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
