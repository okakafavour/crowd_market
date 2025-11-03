package services

import (
	"CROWD_MARKET/config"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateJWT creates a JWT token for the given email
func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyGoogleIdToken validates the Google ID token and returns the email
func VerifyGoogleIdToken(idTokenStr string) (string, error) {
	payload, err := idtoken.Validate(context.Background(), idTokenStr, config.GoogleOauthConfig.ClientID)
	if err != nil {
		return "", err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", errors.New("email not found in token")
	}

	return email, nil
}

// CreateJWT is just a wrapper around GenerateJWT
func CreateJWT(email string) (string, error) {
	return GenerateJWT(email)
}
