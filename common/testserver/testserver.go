package testserver

import (
	"net/http/httptest"
	"sushi-mart/api/analytics"
	"sushi-mart/api/inventory"
	"sushi-mart/api/orders"
	"sushi-mart/api/user"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"sushi-mart/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinServer struct {
	*httptest.Server
	//Config      *common.Config
	//RouterGroup *gin.RouterGroup
}

func GetNewTestGinServer() *GinServer {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middlewares.CORSMiddleware())

	//setup all the routes during tests init

	s := httptest.NewUnstartedServer(r)
	return &GinServer{
		Server: s,
	}
}

func (s *GinServer) SetupRoutes(engine *gin.Engine, queries database.Querier, config *common.Config, logger *logrus.Logger) {
	router := engine.Group("/api/v1")

	//default middlewares
	router.Use(middlewares.LoggerMiddleware(logger))

	//routergroup for managing inventory, restricted to admins
	inventoryRouterGrp := router.Group("/admin/inventory", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	inventory.New(queries).HandleInventory(inventoryRouterGrp)

	//routergroup to check users,orders,products analytics, restricted to admins
	analyticsRouterGrp := router.Group("/admin/analytics", gin.BasicAuth(gin.Accounts{
		config.AdminUser: config.AdminPass,
	}))
	analytics.New(queries).HandleAnalytics(analyticsRouterGrp)

	//routergroup for managing users
	users := router.Group("/users")
	user.New(queries).HandleUsers(users, config)

	//jwt authenticated routes
	router.Use(middlewares.JwtMiddleware(config))
	orderRouterGrp := router.Group("/orders")
	orders.New(queries).HandleOrders(orderRouterGrp, config)
}
