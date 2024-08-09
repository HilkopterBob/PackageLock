package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Used for Agent-SignUps
// var agent_secret_key string = "Secret_Key"

// Data structs

// TODO: support for multiple network adapters.
type Network_Info struct {
	Ip_addr  string
	Mac_addr string
	// TODO: add domain or FQDN
}

type Package_Manager struct {
	Package_manager_name string
	Package_repos        []string
}

type Host struct {
	// TODO: support different linux distros
	ID               int
	Name             string
	Current_packages []string
	Network_info     Network_Info
	Package_manager  Package_Manager
}

type Agent struct {
	Agent_name   string
	Agent_secret string
	Host_ID      int
	Agent_ID     int
}

var hosts = []Host{
	{ID: 1, Name: "Host1", Network_info: Network_Info{Ip_addr: "192.168.1.1", Mac_addr: "AA:BB:CC:DD:EE:FF"}, Package_manager: Package_Manager{Package_manager_name: "pacman", Package_repos: []string{"Repo1", "Repo2", "Repo3"}}, Current_packages: []string{"Package1", "package2", "Package3"}},
	{ID: 2, Name: "Host2", Network_info: Network_Info{Ip_addr: "192.168.1.2", Mac_addr: "AA:BB:CC:DD:EF:00"}, Package_manager: Package_Manager{Package_manager_name: "pacman", Package_repos: []string{"Repo1", "Repo2", "Repo3"}}, Current_packages: []string{"Package1", "package2", "Package3"}},
	{ID: 3, Name: "Ich liebe dich", Network_info: Network_Info{Ip_addr: "192.168.1.3", Mac_addr: "AA:BB:CC:DD:EF:01"}, Package_manager: Package_Manager{Package_manager_name: "pacman", Package_repos: []string{"Repo1", "Repo2", "Repo3"}}, Current_packages: []string{"Package1", "package2", "Package3"}},
}

var agents = []Agent{
	{Agent_name: "Agent Host1", Agent_secret: "11:11:11:11", Host_ID: 1, Agent_ID: 1},
	{Agent_name: "Agent Host2", Agent_secret: "11:11:11:12", Host_ID: 2, Agent_ID: 2},
	{Agent_name: "Agent Host3", Agent_secret: "11:11:11:13", Host_ID: 3, Agent_ID: 3},
}

func getHosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, hosts)
}

func getAgents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, agents)
}

func getAgentByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range agents {
		if strconv.Itoa(a.Host_ID) == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no agent under that id"})
}

// POST Functions
func registerAgent(c *gin.Context) {
	var newAgent Agent

	if err := c.BindJSON(&newAgent); err != nil {
		// TODO: Add logs
		// TODO: Add errorhandling
		return
	}

	agents = append(agents, newAgent)
	c.IndentedJSON(http.StatusCreated, newAgent)
}

func registerHost(c *gin.Context) {
	var newHost Host

	if err := c.BindJSON(&newHost); err != nil {
		// TODO: Add logs
		// TODO: Add errorhandling
		return
	}

	hosts = append(hosts, newHost)
	c.IndentedJSON(http.StatusCreated, newHost)
}

func getHostByAgentID(c *gin.Context) {
	var agent_by_id Agent

	// gets the value from /agent/:id/host
	id := c.Param("id")

	// finds the agent by the URL-ID
	for _, a := range agents {
		if strconv.Itoa(a.Host_ID) == id {
			// c.IndentedJSON(http.StatusOK, a)
			agent_by_id = a
		}
	}

	// finds host with same id as agent
	for _, host := range hosts {
		if host.ID == agent_by_id.Agent_ID {
			c.IndentedJSON(http.StatusOK, host)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no agent under that id"})
}

func main() {
	// Endpoints & Data Aggregation Functions

	//  API v0.1 structure:
	//  /hosts
	//  GET: ✅
	//    - shows all hosts and the hosts data
	//  POST: ✅
	//    - adds new host to 'hosts'-slice
	//
	//
	//  /agents
	//  GET: ✅
	//    - shows all agents and the agents data
	//  POST: ✅
	//    - adds new agent to 'agents'-slice
	//  /agent/:id/host ✅
	//  GET:
	//    - shows the host connected to the agent
	//  /agent/:id
	//  GET: ✅
	//    - shows agent with
	//
	//  /commandqueue/agent
	//  GET:
	//    - respond with 'no commands' or 'new commands'
	//  POST:
	//    - post Agent.agent_secret_key, respond with commands

	router := gin.Default()
	router.GET("/hosts", getHosts)
	router.POST("/hosts", registerHost)
	router.GET("/agents", getAgents)
	router.POST("/agents", registerAgent)
	router.GET("/agent/:id", getAgentByID)
	router.GET("/agent/:id/host", getHostByAgentID)

	// TODO: create logs
	// TODO: write error to logs
	// TODO: handle error 'This port is blocked, check your FW or smth'
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
}
