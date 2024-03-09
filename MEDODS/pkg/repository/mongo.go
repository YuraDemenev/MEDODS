package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const usersTable = "users"

type ConfigDB struct {
	Url    string
	Addres string
}

type User struct {
	GUID             string
	Password         string
	UserName         string
	RefreshTokenGUID string
	RefreshToken     []byte
	ExpiresAt        time.Time
}

func NewMongoDB(configDB ConfigDB, ctx context.Context) (*mongo.Client, error) {
	uri := fmt.Sprintf("%s://%s", configDB.Url, configDB.Addres)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil

}
