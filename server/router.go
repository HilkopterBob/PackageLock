package server

import (
	"packagelock/handler"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Router *gin.Engine
}

func (router Routes) addAgentHandlers(rg *gin.RouterGroup) {
	AgentGroup := rg.Group("/agent")

	AgentGroup.GET("/:id", handler.GetAgentByID)
	AgentGroup.GET("/:id/host", handler.GetHostByAgentID)
	AgentGroup.POST("/register", handler.RegisterAgent)
}

func (router Routes) addGeneralHandlers(rg *gin.RouterGroup) {
	GeneralGroup := rg.Group("/general")

	GeneralGroup.GET("/hosts", handler.GetHosts)
	GeneralGroup.GET("/agents", handler.GetAgents)
}

func (router Routes) addHostHandlers(rg *gin.RouterGroup) {
	HostGroup := rg.Group("/host")

	HostGroup.POST("/register", handler.RegisterHost)
}

func AddRoutes() Routes {
	router := Routes{
		Router: gin.Default(),
	}

	v1 := router.Router.Group("/v1")

	router.addGeneralHandlers(v1)
	router.addAgentHandlers(v1)
	router.addHostHandlers(v1)

	return router
}
