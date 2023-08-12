package orders

import "github.com/gin-gonic/gin"

func (wrapper *RoutesWrapper) HandleOrders(router *gin.RouterGroup) {
	router.POST("/place-order", wrapper.HandlePlaceOrder)
	router.POST("/cancel-order", wrapper.HandleCancelOrder)
	router.GET("/get-orders", wrapper.HandleGetOrders)
}
