package server

import (
	"packagelock/handler"

	"github.com/gin-gonic/gin"
)

// 'Routes' holds the gin-Engine Pointer.
type Routes struct {
	Router *gin.Engine
}

func (router Routes) addAgentHandler(rg *gin.RouterGroup) {
	AgentGroup := rg.Group("/agent")

	AgentGroup.GET("/:id", handler.GetAgentByID)
	AgentGroup.GET("/:id/host", handler.GetHostByAgentID)
	AgentGroup.POST("/register", handler.RegisterAgent)
}

func (router Routes) addGeneralHandler(rg *gin.RouterGroup) {
	GeneralGroup := rg.Group("/general")

	GeneralGroup.GET("/hosts", handler.GetHosts)
	GeneralGroup.GET("/agents", handler.GetAgents)
}

func (router Routes) addHostHandler(rg *gin.RouterGroup) {
	HostGroup := rg.Group("/host")

	HostGroup.POST("/register", handler.RegisterHost)
}

// AddRoutes adds all handler groups to the current router.
// Its Exported, used in main() and returns a Router typed Routes.
// AddRoutes calls all add_handlergroupname_Handler functions.
func AddRoutes() Routes {
	router := Routes{
		Router: gin.Default(),
	}

	v1 := router.Router.Group("/v1")

	router.addGeneralHandler(v1)
	router.addAgentHandler(v1)
	router.addHostHandler(v1)

	return router
}
