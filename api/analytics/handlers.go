package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *RoutesWrapper) HandleAvgRatings(c *gin.Context) {
	resp, err := r.AnalyticsService.GetAvgCustomerRatings(c.Request.Context())
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (r *RoutesWrapper) HandleOrdersPlaced(c *gin.Context) {
	limitStr, ok := c.GetQuery("limit")
	if !ok {
		c.JSON(http.StatusBadRequest, "limit param query required")
		return
	}
	limit, er := strconv.Atoi(limitStr)
	if er != nil {
		c.JSON(http.StatusBadRequest, "limit param query required")
		return
	}

	resp, err := r.AnalyticsService.GetMostOrdersPlaced(c.Request.Context(), limit)
	if err != nil {
		c.JSON(err.Status, err.Message)
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}
