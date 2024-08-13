package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"packagelock/structs"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func Test_registerAgent(t *testing.T) {
	router := SetUpRouter()
	router.POST("/v1/agent/register", registerAgent)
	test_agent := structs.Agent{
		Agent_name:   "Test Agent",
		Agent_secret: "FF:FF:FF:FF:FF:FF",
		Host_ID:      9,
		Agent_ID:     99,
	}

	json_test_agent, _ := json.Marshal(test_agent)
	request, _ := http.NewRequest("POST", "/v1/agent/register", bytes.NewBuffer(json_test_agent))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusCreated, response.Code)
}

func Test_registerHost(t *testing.T) {
	// FIXME: This test passes with broken host objects!
	// FIXME: Control/COrrect Other tests!

	router := SetUpRouter()
	router.POST("/v1/host/register", registerHost)
	test_Host := structs.Host{
		ID:   99,
		Name: "Testhost",
		Network_info: structs.Network_Info{
			Ip_addr:  "192.168.1.3",
			Mac_addr: "AA:BB:CC:DD:EF:01",
		},
		Package_manager: structs.Package_Manager{
			Package_manager_name: "pacman",
			Package_repos: []string{
				"Repo1",
				"Repo2",
				"Repo3",
			},
		},
		Current_packages: []string{
			"Package1",
			"package2",
			"Package3",
		},
	}

	json_test_host, _ := json.Marshal(test_Host)
	request, _ := http.NewRequest("POST", "/v1/host/register", bytes.NewBuffer(json_test_host))

	fmt.Println(router.GET("/v1/general/hosts", getAgents))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusCreated, response.Code)
}

func Test_getAgents(t *testing.T) {
	router := SetUpRouter()
	router.GET("/v1/general/agents", getAgents)
	request, _ := http.NewRequest("GET", "/v1/general/agents", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var agents []structs.Agent
	err := json.Unmarshal(response.Body.Bytes(), &agents)
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotEmpty(t, agents)
}

func Test_getHosts(t *testing.T) {
	router := SetUpRouter()
	router.GET("/v1/general/hosts", getHosts)
	request, _ := http.NewRequest("GET", "/v1/general/hosts", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var hosts []structs.Host
	err := json.Unmarshal(response.Body.Bytes(), &hosts)
	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotEmpty(t, hosts)
}

func Test_getAgentByID(t *testing.T) {
	router := SetUpRouter()
	router.GET("/v1/agent/:id", getAgentByID)
	request, _ := http.NewRequest("GET", "/v1/agent/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var agent structs.Agent
	err := json.Unmarshal(response.Body.Bytes(), &agent)
	if err != nil {
		fmt.Println(response.Body)
		fmt.Println(&agent)
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotEmpty(t, agent)
}

func Test_getHostByAgentID(t *testing.T) {
	router := SetUpRouter()
	router.GET("/v1/agent/:id/host", getHostByAgentID)
	request, _ := http.NewRequest("GET", "/v1/agent/1/host", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var host structs.Host
	err := json.Unmarshal(response.Body.Bytes(), &host)
	if err != nil {
		fmt.Println(response.Body)
		fmt.Println(host)
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotEmpty(t, host)
}
