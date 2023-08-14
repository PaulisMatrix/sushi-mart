package orders

import (
	"sushi-mart/common"

	"github.com/gin-gonic/gin"
)

func (wrapper *RoutesWrapper) HandleOrders(router *gin.RouterGroup, config *common.Config) {
	router.POST("/place-order", wrapper.HandlePlaceOrder(config))
	router.POST("/cancel-order", wrapper.HandleCancelOrder)
	router.GET("/get-orders", wrapper.HandleGetOrders)
}
