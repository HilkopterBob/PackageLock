package handler

import (
	"net/http"
	"packagelock/test_data"

	"github.com/gin-gonic/gin"
)

func GetHosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, test_data.Hosts)
}

func GetAgents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, test_data.Agents)
}
