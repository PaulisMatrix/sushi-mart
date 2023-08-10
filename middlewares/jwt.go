package middlewares

import (
	"fmt"
	"net/http"
	"sushi-mart/common"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtMiddleware(config *common.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token != "" {

			//validate the jwt token
			jwtToken, err := parseToken(token, config)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
				return
			}

			_, err = GetClaims(jwtToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
			return
		}

		c.Next()
	}
}

// ParseToken parse a token
func parseToken(tokenString string, config *common.Config) (*jwt.Token, error) {
	var token *jwt.Token
	var err error
	token, err = parseHS256(tokenString, token, config)
	if err != nil {
		return nil, err
	}

	return token, err
}

// GetClaims get claims information
func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}
	err := token.Claims.(jwt.MapClaims).Valid()
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func parseHS256(tokenString string, token *jwt.Token, config *common.Config) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSktKey), nil
	})
	return token, err
}
