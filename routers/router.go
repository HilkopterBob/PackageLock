package routers

import "github.com/gin-gonic/gin"

type Routes struct {
	Router *gin.Engine
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
