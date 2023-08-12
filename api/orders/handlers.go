package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Place an Order
// @Description Places a Customer Order
// @Schemes http
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body PlaceOrderReq true "PlaceOrderRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/place-order [post]
func (r *RoutesWrapper) HandlePlaceOrder(c *gin.Context) {
	//get userID from gin context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "userID missing in the context"})
		return
	}

	custId, isok := userID.(int)
	if !isok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "userID not of type int"})
		return
	}

	var input PlaceOrderReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	err := r.OrderService.PlaceOrder(c.Request.Context(), &input, custId)

	if err != nil {
		c.JSON(err.Status, gin.H{"message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order placed successfully"})
	return
}

// @Summary Cancel an Order
// @Description Cancels a Customer Order
// @Schemes http
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body UpdateOrderReq true "CancelOrderRequest"
// @Success 200 {string} SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/cancel-order [post]
func (r *RoutesWrapper) HandleCancelOrder(c *gin.Context) {
	var input UpdateOrderReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	err := r.OrderService.CancelOrder(c.Request.Context(), &input)

	if err != nil {
		c.JSON(err.Status, gin.H{"message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
	return
}

// @Summary Get Orders
// @Description Returns all Orders placed by the Customer
// @Schemes http
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} GetAllOrdersResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/get-orders [get]
func (r *RoutesWrapper) HandleGetOrders(c *gin.Context) {
	//get userID from gin context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "userID missing in the context"})
		return
	}

	custId, isok := userID.(int)
	if !isok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "userID not of type int"})
		return
	}

	resp, err := r.OrderService.GetOrders(c.Request.Context(), custId)

	if err != nil {
		c.JSON(err.Status, gin.H{"message": err.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
}
