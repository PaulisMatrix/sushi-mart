package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Returns Average Customer Ratings
// @Description Returns average customer ratings for the orders pucharsed by them
// @Schemes http
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 200 {object} AvgCustomerRatingsResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/analytics/avg-cust-ratings [get]
func (r *RoutesWrapper) HandleAvgRatings(c *gin.Context) {
	resp, err := r.AnalyticsService.GetAvgCustomerRatings(c.Request.Context())
	if err != nil {
		c.JSON(err.Status, gin.H{"message": err.Message})
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// @Summary Returns Customers Placed Orders
// @Description Returns the most placed orders by the customers
// @Schemes http
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param limit query int true "limit"
// @Success 200 {object} MostOrdersPlacedResp
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 429 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /admin/analytics/top-orders-placed [get]
func (r *RoutesWrapper) HandleOrdersPlaced(c *gin.Context) {
	limitStr, ok := c.GetQuery("limit")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "limit param query required"})
		return
	}
	limit, er := strconv.Atoi(limitStr)
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "limit param query required"})
		return
	}

	resp, err := r.AnalyticsService.GetMostOrdersPlaced(c.Request.Context(), limit)
	if err != nil {
		c.JSON(err.Status, gin.H{"message": err.Message})
		return
	}

	if resp.OrdersPlaced == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not a single order placed by a customer yet"})
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}
