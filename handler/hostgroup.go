package handler

import (
	"net/http"
	"packagelock/structs"

	"github.com/gin-gonic/gin"
)

func RegisterHost(c *gin.Context) {
	var newHost structs.Host

	if err := c.BindJSON(&newHost); err != nil {
		// TODO: Add logs
		// TODO: Add errorhandling
		return
	}

	Hosts = append(Hosts, newHost)
	c.IndentedJSON(http.StatusCreated, newHost)
}
