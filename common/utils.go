package common

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
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

type LoggerKey struct{}

// extract logger unsafe
func ExtractLoggerUnsafe(ctx context.Context) *logrus.Logger {
	switch ctx.Value(LoggerKey{}).(type) {
	case *logrus.Logger:
		lg, ok := ctx.Value(LoggerKey{}).(*logrus.Logger)
		if !ok {
			return logrus.StandardLogger()
		}
		return lg
	default:
		return logrus.StandardLogger()
	}
}
