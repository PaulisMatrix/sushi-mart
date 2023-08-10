package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func (r *RoutesWrapper) SignUp(c *gin.Context) {
	var custInput SignUpReq

	if err := c.ShouldBindJSON(&custInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(custInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := database.CreateCustomerParams{
		Username: custInput.Username,
		Password: string(hashedPassword),
		Email:    custInput.Email,
		Phone:    sql.NullString{String: custInput.Phone, Valid: true},
		Address:  sql.NullString{String: custInput.Address, Valid: true},
	}

	resp, err := r.UsersService.CreateUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user added to the db with username %s", resp.Username)})
	return
}

func (r *RoutesWrapper) Login(config *common.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var custInput LoginReq

		if err := c.ShouldBindJSON(&custInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := r.UsersService.GetUser(c.Request.Context(), custInput.Email)
		if err != nil && err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "signup first"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(custInput.Password))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials!!"})
			return
		}

		//generate 1hr long token and return
		token, err := generateNewToken(resp.Email, config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		loginRes := &LoginResp{
			Username: resp.Username,
			JWTToken: token,
		}
		c.JSON(http.StatusOK, loginRes)
		return
	}
}

func generateNewToken(userId string, config *common.Config) (string, error) {
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
