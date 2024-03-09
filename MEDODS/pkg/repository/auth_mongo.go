package repository

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthMongo struct {
	mongo *mongo.Client
}

func NewAuthMongo(mongo *mongo.Client) *AuthMongo {
	return &AuthMongo{mongo: mongo}
}

// Create new User
func (a *AuthMongo) CreateUser(userName, password string) (string, error) {
	//If our GUID Exist
	check := true
	guid := ""
	for check {
		guid = a.generateGUID()
		check = a.checkGuidExist(guid)
	}

	user := User{GUID: guid, Password: password, UserName: userName,
		RefreshTokenGUID: "", RefreshToken: []byte{}, ExpiresAt: time.Now().UTC()}

	//Insert User
	collection := a.mongo.Database("usersDB").Collection(usersTable)
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}

	return "", nil
}

// Get user id
func (a *AuthMongo) GetUserId(userName, password string) (string, error) {
	collection := a.mongo.Database("usersDB").Collection(usersTable)
	filter := bson.D{{Key: "username", Value: userName}, {Key: "password", Value: password}}

	var result User

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", nil
	}

	return result.GUID, nil
}

// Generete GUID
func (a *AuthMongo) generateGUID() string {
	guid := uuid.New()
	return guid.String()
}

// Check Does GUID exist in db
func (a *AuthMongo) checkGuidExist(guid string) bool {
	//Connect to db
	ctx := context.Background()
	collection := a.mongo.Database("usersDB").Collection(usersTable)

	filter := bson.M{"guid": guid}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logrus.Fatalf("Cant check guid exist, Err:%s", err.Error())
		return true
	}

	if count == 0 {
		return false
	} else {
		return true
	}

}

// Save Refresh token
func (a *AuthMongo) SaveRefreshToken(refreshToken string, guid string) (string, error) {
	//Connect to db
	collection := a.mongo.Database("usersDB").Collection(usersTable)
	filter := bson.D{{Key: "guid", Value: guid}}

	//Decode because bcrypt dont get >72 bytes
	tokenBytes, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return "", err
	}

	//crypt to bytes for save refresh token
	bytes, err := bcrypt.GenerateFromPassword(tokenBytes, 10)
	if err != nil {
		return "", err
	}

	//Create GUID for refresh token
	refreshTokenGUID := a.generateGUID()

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "refreshtoken", Value: bytes},
		{Key: "expiresat", Value: time.Now().Local().Add(time.Hour * 24 * 7).UTC()},
		{Key: "refreshtokenguid", Value: refreshTokenGUID}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return "", err
	}

	return refreshTokenGUID, nil
}

func (a *AuthMongo) GetUserName(guid string) (string, error) {
	//Connect to db
	collection := a.mongo.Database("usersDB").Collection(usersTable)
	filter := bson.D{{Key: "guid", Value: guid}}

	var result User

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", nil
	}

	return result.UserName, nil
}
