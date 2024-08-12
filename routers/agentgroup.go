package routers

import (
	"packagelock/handler"

	"github.com/gin-gonic/gin"
)

func (router Routes) addAgentHandlers(rg *gin.RouterGroup) {
	AgentGroup := rg.Group("/agent")

	AgentGroup.GET("/:id", handler.GetAgentByID)
	AgentGroup.GET("/:id/host", handler.GetHostByAgentID)
	AgentGroup.POST("/register", handler.RegisterAgent)
}
