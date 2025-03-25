package util

import (
	"database/sql"
	"errors"
	"fmt"
	"gin-api/database"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid token")
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func loadSecretKey() []byte {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	return []byte(secretKey)
}

func GenerateTokenPair(userId int, isAdmin bool, isEgnLab bool) (*TokenPair, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["userId"] = userId
	accessClaims["isAdmin"] = isAdmin
	accessClaims["isEgnLab"] = isEgnLab
	accessClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	accessClaims["type"] = "access"

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["userId"] = userId
	refreshClaims["isAdmin"] = isAdmin
	refreshClaims["isEgnLab"] = isEgnLab
	refreshClaims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshClaims["type"] = "refresh"

	jwtSecret := loadSecretKey()

	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		jwtSecret := loadSecretKey()
		return jwtSecret, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	return token, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	// Validate refresh token
	token, err := ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	// Verify it's a refresh token
	if claims["type"] != "refresh" {
		return "", ErrInvalidToken
	}

	// Generate new access token
	userId := int(claims["userId"].(float64))
	// Extract userId from claims

	// Query the database to get latest user data
	var isAdmin, isEgnLab bool
	err = database.DB.QueryRow("SELECT admin, is_egn_lab FROM users WHERE id = ?", userId).Scan(
		&isAdmin,
		&isEgnLab,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user no longer exists")
		}
		return "", fmt.Errorf("database error: %v", err)
	}

	newToken := jwt.New(jwt.SigningMethodHS256)
	newClaims := newToken.Claims.(jwt.MapClaims)
	newClaims["userId"] = userId
	newClaims["isAdmin"] = isAdmin
	newClaims["isEgnLab"] = isEgnLab
	newClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	newClaims["type"] = "access"

	jwtSecret := loadSecretKey()
	return newToken.SignedString(jwtSecret)
}
