package services

import (
	"context"
	"errors"
	"time"

	"CROWD_MARKET/config"
	"CROWD_MARKET/model"
	"CROWD_MARKET/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection // global variable

func InitUserService() {
	userCollection = config.DB.Collection("users")
}

func RegisterUser(name, email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return err
	}

	newUser := model.User{
		ID:               primitive.NewObjectID(),
		Name:             name,
		Email:            email,
		Password:         string(hashedPassword),
		IsVerified:       false,
		VerificationCode: code,
		Provider:         "local",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return err
	}

	return SendVerificationEmail(email, code)
}

func VerifyUserEmail(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"verificationCode": code}
	update := bson.M{
		"$set": bson.M{
			"isVerified":       true,
			"verificationCode": "",
			"updatedAt":        time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("invalid or expired verification code")
	}
	return nil
}

func LoginUser(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	if !user.IsVerified {
		return "", errors.New("email is not verified")
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func FindUserByEmail(email string) (*model.User, error) {
	var user model.User

	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateGoogleUser(name, email string) error {
	user := model.User{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Email:      email,
		IsVerified: true,
		Provider:   "google",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := userCollection.InsertOne(context.Background(), user)
	return err
}
