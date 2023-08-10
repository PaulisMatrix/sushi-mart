package user

import "github.com/golang-jwt/jwt/v4"

type SignUpReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

type LoginReq struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LoginResp struct {
	Username string `json:"username"`
	JWTToken string `json:"token"`
}

type CustomClaims struct {
	UserId string `json:"userID"`
	jwt.StandardClaims
}
