package user

import (
	"net/http"
	"sushi-mart/common"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary User Signup
// @Description Register a Customer
// @Schemes http
// @Accept json
// @Produce json
// @Param data body SignUpReq true "UserSignupRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/signup [post]
func (r *RoutesWrapper) SignUp(c *gin.Context) {
	var input SignUpReq

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := r.UsersService.CreateUser(c.Request.Context(), &input)
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, "new user created successfully")
	return
}

// @Summary User Login
// @Description Login a Customer and generate a new JWT Token
// @Schemes http
// @Accept json
// @Produce json
// @Param data body LoginReq true "UserLoginRequest"
// @Success 200 {object} LoginResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/login [post]
func (r *RoutesWrapper) Login(config *common.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginReq
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		resp, err := r.UsersService.GetUser(c.Request.Context(), &input)
		if err != nil {
			c.JSON(err.Status, err.Message)
			return
		}

		hashErr := bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(input.Password))
		if hashErr != nil && hashErr == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusForbidden, "invalid credentials")
			return
		}

		//generate 1hr long token and return
		token, tokenErr := common.GenerateNewToken(int32(resp.CustId), config)
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, "internal server error")
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

// @Summary Create Customer Wallet
// @Description Create a new Wallet attached to the Customer
// @Schemes http
// @Accept json
// @Produce json
// @Param data body CreateWalletReq true "CreateWalletRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/create-wallet [post]
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

// @Summary Get Wallet
// @Description Returns Wallet attached to the Customer
// @Schemes http
// @Accept json
// @Produce json
// @Success 200 {object} GetWalletRes
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/get-wallet [get]
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

// @Summary Update Customer Wallet
// @Description Update different fields like Balance, PaymentType, etc of a Customer
// @Schemes http
// @Accept json
// @Produce json
// @Param data body UpdateWalletReq true "UpdateWalletRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/update-wallet [patch]
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

// @Summary Get ProductItems
// @Description Returns a list of all ProductItems for Customer to select from
// @Schemes http
// @Accept json
// @Produce json
// @Success 200 {object} GetAllProductsResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/all-products [get]
func (r *RoutesWrapper) HandleAllProducts(c *gin.Context) {
	resp, err := r.UsersService.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// @Summary Add a Review
// @Description Registers Cutomers Reviews for a particular ProductItem
// @Schemes http
// @Accept json
// @Produce json
// @Param data body AddReviewReq true "AddReviewRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /users/add-review [post]
func (r *RoutesWrapper) HandleAddReview(c *gin.Context) {
	//get userID from gin context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, "userID missing in the context")
		return
	}

	custId, isok := userID.(int)
	if !isok {
		c.JSON(http.StatusBadRequest, "userID not of type int")
		return
	}

	var input AddReviewReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	err := r.UsersService.AddReview(c.Request.Context(), &input, custId)

	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, "successfully added your review")
	return
}
