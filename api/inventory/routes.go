package inventory

import (
	"github.com/gin-gonic/gin"
)

func (wrapper *RoutesWrapper) HandleInventory(router *gin.RouterGroup) {
	router.GET("/all-products", wrapper.HandleAllProducts)
	router.POST("/add-product", wrapper.HandleAddProduct)
	router.PATCH("/update-product/:id", wrapper.HandleUpdateProduct)
	router.DELETE("/delete-product/:id", wrapper.HandleDeleteProduct)
}
