package user

import (
	"sushi-mart/common"

	"github.com/gin-gonic/gin"
)

func (wrapper *RoutesWrapper) HandlUsers(router *gin.RouterGroup, config *common.Config) {
	router.POST("/signup", wrapper.SignUp)
	router.POST("/login", wrapper.Login(config))
}
