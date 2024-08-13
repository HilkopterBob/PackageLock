package handler

import (
	"net/http"
	"packagelock/structs"
	"packagelock/test_data"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAgentByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range test_data.Agents {
		if strconv.Itoa(a.Host_ID) == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no agent under that id"})
}

// POST Functions
func RegisterAgent(c *gin.Context) {
	var newAgent structs.Agent

	if err := c.BindJSON(&newAgent); err != nil {
		// TODO: Add logs
		// TODO: Add errorhandling
		return
	}

	test_data.Agents = append(test_data.Agents, newAgent)
	c.IndentedJSON(http.StatusCreated, newAgent)
}

func GetHostByAgentID(c *gin.Context) {
	var agent_by_id structs.Agent

	// gets the value from /agent/:id/host
	id := c.Param("id")

	// finds the agent by the URL-ID
	for _, a := range test_data.Agents {
		if strconv.Itoa(a.Host_ID) == id {
			// c.IndentedJSON(http.StatusOK, a)
			agent_by_id = a
		}
	}

	// finds host with same id as agent
	for _, host := range test_data.Hosts {
		if host.ID == agent_by_id.Agent_ID {
			c.IndentedJSON(http.StatusOK, host)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no agent under that id"})
}