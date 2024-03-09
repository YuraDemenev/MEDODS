package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type RefreshMongo struct {
	mongo *mongo.Client
}

func NewRefreshMongo(mongo *mongo.Client) *RefreshMongo {
	return &RefreshMongo{mongo: mongo}
}

func (r *RefreshMongo) CheckExpireTime(tokenRefresh string, tokenRefreshGUID string) (string, error) {
	collection := r.mongo.Database("usersDB").Collection(usersTable)
	fmt.Print("Create refresh token")

	//Decode because bcrypt dont get >72 bytes
	tokenBytes, err := base64.StdEncoding.DecodeString(tokenRefresh)
	if err != nil {
		return "", err
	}

	//Create filter and user
	filter := bson.D{{Key: "refreshtokenguid", Value: tokenRefreshGUID}}
	var result User

	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(result.RefreshToken, tokenBytes)
	if err != nil {
		logrus.Error("Cant compare hash and password bytes")
		return "", err
	}

	//Check expire time
	today := time.Now().UTC()
	check := today.Before(result.ExpiresAt)

	if check {
		return result.GUID, nil
	} else {
		return "", errors.New("time expire")
	}

}
