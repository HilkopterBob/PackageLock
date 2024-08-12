package routers

import (
	"packagelock/handler"

	"github.com/gin-gonic/gin"
)

func (router Routes) addHostHandlers(rg *gin.RouterGroup) {
	HostGroup := rg.Group("/host")

	HostGroup.POST("/register", handler.RegisterHost)
}
