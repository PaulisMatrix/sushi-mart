package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"

	"github.com/gin-gonic/gin"
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
		token, err := common.GenerateNewToken(resp.ID, config)
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

func (r *RoutesWrapper) HandleCreateWallet(c *gin.Context) {
	//get userID from gin context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, "userID missing in the context")
		return
	}

	Id, isok := userID.(int)
	if !isok {
		c.JSON(http.StatusBadRequest, "userID not of type int")
		return
	}

	var input CreateWalletReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := r.UsersService.CreateUserWallet(c.Request.Context(), &input, Id)

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, "succesfully created the wallet")
	return

}

func (r *RoutesWrapper) HandleGetWallet(c *gin.Context) {
	//get userID from gin context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, "userID missing in the context")
		return
	}

	Id, isok := userID.(int)
	if !isok {
		c.JSON(http.StatusBadRequest, "userID not of type int")
		return
	}

	resp, err := r.UsersService.GetUserWallet(c.Request.Context(), Id)

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (r *RoutesWrapper) HandleUpdateWallet(c *gin.Context) {
	//get userID from gin context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, "userID missing in the context")
		return
	}

	Id, isok := userID.(int)
	if !isok {
		c.JSON(http.StatusBadRequest, "userID not of type int")
		return
	}

	var input UpdateWalletReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := r.UsersService.UpdateUserWallet(c.Request.Context(), &input, Id)

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, "successfully updated your wallet")
	return
}

func (r *RoutesWrapper) HandleAllProducts(c *gin.Context) {
	resp, err := r.UsersService.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}
