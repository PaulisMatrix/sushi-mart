package analytics

import "github.com/gin-gonic/gin"

func (wrapper *RoutesWrapper) HandleAnalytics(router *gin.RouterGroup) {
	router.GET("/avg-cust-ratings", wrapper.HandleAvgRatings)
	router.GET("/top-orders-placed", wrapper.HandleOrdersPlaced)
}
