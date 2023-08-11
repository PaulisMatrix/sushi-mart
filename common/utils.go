package common

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CustomClaims struct {
	UserId int32 `json:"userID"`
	jwt.StandardClaims
}

func GenerateNewToken(userId int32, config *Config) (string, error) {
	customClaim := CustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)
	sToken, err := token.SignedString([]byte(config.JwtSktKey))
	if err != nil {
		return "", err
	}
	return sToken, nil
}
