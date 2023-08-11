package user

import (
	"sushi-mart/common"

	"github.com/gin-gonic/gin"
)

func (wrapper *RoutesWrapper) HandleUsers(router *gin.RouterGroup, config *common.Config) {
	router.POST("/signup", wrapper.SignUp)
	router.POST("/login", wrapper.Login(config))

	//for creating wallet/adding balance
	router.POST("/create-wallet", wrapper.HandleCreateWallet)
	router.GET("/get-wallet", wrapper.HandleGetWallet)
	router.PATCH("/update-wallet", wrapper.HandleUpdateWallet)

	//to select different products
	router.GET("/all-products", wrapper.HandleAllProducts)
}
