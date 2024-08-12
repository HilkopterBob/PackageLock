package routers

import (
	"packagelock/handler"

	"github.com/gin-gonic/gin"
)

func (router Routes) addGeneralHandlers(rg *gin.RouterGroup) {
	GeneralGroup := rg.Group("/general")

	GeneralGroup.GET("/hosts", handler.GetHosts)
	GeneralGroup.GET("/agents", handler.GetAgents)
}
