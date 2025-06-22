package shared

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID uuid.UUID
	Email  string
}

func GenerateAccessToken(tokenClaims *TokenClaims) (string, error) {
	env, err := ParseEnv()
	if err != nil {
		return "", err
	}

	return generateJWT(tokenClaims, env.AccessTokenSecret)
}

func VerifyAccessToken(tokenString string) (*TokenClaims, error) {
	env, err := ParseEnv()
	if err != nil {
		return nil, err
	}

	return verifyJWT(tokenString, env.AccessTokenSecret)
}

func verifyJWT(tokenString string, secretKey string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("Error parsing token: ", err)
		return nil, errors.New(ErrInvalidToken)
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return nil, errors.New(ErrInvalidToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Error parsing claims from token")
		return nil, errors.New(ErrInvalidToken)
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		log.Println("Error parsing userId from token")
		return nil, errors.New(ErrInvalidToken)
	}

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("Error parsing email from token")
		return nil, errors.New(ErrInvalidToken)
	}

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		log.Println("Error parsing userId to UUID: ", err)
		return nil, errors.New(ErrInvalidToken)
	}

	return &TokenClaims{
		UserID: userIdUUID,
		Email:  email,
	}, nil
}

func generateJWT(tokenClaims *TokenClaims, secretKey string) (string, error) {
	env, err := ParseEnv()
	if err != nil {
		return "", err
	}

	expDuration := time.Second * time.Duration(env.AccessTokenExpSeconds)

	claims := jwt.MapClaims{
		"userId": tokenClaims.UserID.String(),
		"email":  tokenClaims.Email,
		"exp":    time.Now().Add(expDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
