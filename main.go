package main

import (
	"fmt"
	"packagelock/config"
	"packagelock/routers"

	"github.com/spf13/viper"
)

// Data structs

// TODO: support for multiple network adapters.

func main() {
	config.StartViper()
	fmt.Println(viper.AllSettings())

	router := routers.AddRoutes()
	err := router.Router.Run("localhost:9090")
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Endpoints & Data Aggregation Functions

	//  API v0.1 structure,:
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

	// TODO: Group Routes via:
	// https://stackoverflow.com/questions/62906766/how-to-group-routes-in-gin
	// router := gin.Default()
	// router.GET("/hosts", getHosts)
	// router.POST("/hosts", registerHost)
	// router.GET("/agents", getAgents)
	// router.POST("/agents", registerAgent)
	// router.GET("/agent/:id", getAgentByID)
	// router.GET("/agent/:id/host", getHostByAgentID)

	// TODO: create logs
	// TODO: write error to logs
	// TODO: handle error 'This port is blocked, check your FW or smth'

	// TODO: use FQDN and Port from config file
	// fmt.Println(viper.Get("network.fqdn"))
	// err := router.Run(viper.GetString("network.fqdn") + ":" + viper.GetString("network.port"))
	//if err != nil {
	//	fmt.Println(err)
	//}
}
