package middlewares

import (
	"context"
	"sushi-mart/common"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//add logrus logger to gin context

		//get parent context
		parentCtx := context.WithValue(c.Request.Context(), common.LoggerKey{}, logger)

		if parentCtx == nil {
			panic("context logger not being set properly!!")
		}

		//update gin's context
		c.Request = c.Request.WithContext(parentCtx)
		c.Next()
	}
}
