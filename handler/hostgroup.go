package handler

import (
	"net/http"
	"packagelock/structs"
	"packagelock/test_data"

	"github.com/gin-gonic/gin"
)

func RegisterHost(c *gin.Context) {
	var newHost structs.Host

	if err := c.BindJSON(&newHost); err != nil {
		// TODO: Add logs
		// TODO: Add errorhandling
		return
	}

	test_data.Hosts = append(test_data.Hosts, newHost)
	c.IndentedJSON(http.StatusCreated, newHost)
}
