package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Hosts)
}

func GetAgents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Agents)
}
