package handler

import (
	"encoding/base64"
	"packagelock/db"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

// GetAgentByID filters a slice of Agents for a matching Agent.Agent_ID.
// It returns a JSON response with fiber.StatusOK or fiber.StatusNotFound.
func GetAgentByID(c *fiber.Ctx) error {
	// ID is an URL slice. Its a URL-Save base64 encoded UUID
	urlIDBytes, err := base64.RawURLEncoding.DecodeString(c.Params("id"))
	if err != nil {
		// FIXME: error handling
		panic(err)
	}

	urlIDString := string(urlIDBytes)

	agents, err := db.DB.Select("agents")
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}

	var agentsSlice []structs.Agent
	err = surrealdb.Unmarshal(agents, &agentsSlice)
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}

	var requestedAgentByID structs.Agent
	for _, agent := range agentsSlice {
		if agent.AgentID.String() == urlIDString {
			requestedAgentByID = agent
			// FIXME: logging
			return c.Status(fiber.StatusOK).JSON(requestedAgentByID)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(nil)
	// FIXME: logging
}

// RegisterAgent handles POST requests to register a new agent.
func RegisterAgent(c *fiber.Ctx) error {
	var newAgent structs.Agent

	// Parse the JSON request body into newAgent
	if err := c.BodyParser(&newAgent); err != nil {
		// TODO: Add logs
		// TODO: Add error handling
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	newAgentInsertionData, err := db.DB.Create("agents", newAgent)
	if err != nil {
		// FIXME: logging
		// FIXME: error handling
		panic(err)
	}
	// Respond with the newly created agent
	return c.Status(fiber.StatusCreated).JSON(newAgentInsertionData)
}

// GetHostByAgentID finds the host for a given agent ID.
func GetHostByAgentID(c *fiber.Ctx) error {
	urlIDBytes, err := base64.RawURLEncoding.DecodeString(c.Params("id"))
	if err != nil {
		// FIXME: error handling
		panic(err)
	}

	urlIDString := string(urlIDBytes)

	agents, err := db.DB.Select("agents")
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}

	var agentsSlice []structs.Agent
	err = surrealdb.Unmarshal(agents, &agentsSlice)
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}

	var requestedAgentByID structs.Agent
	for _, agent := range agentsSlice {
		if agent.AgentID.String() == urlIDString {
			requestedAgentByID = agent
		}
	}

	if requestedAgentByID.ID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Agent not found",
		})
	}

	hosts, err := db.DB.Select("hosts")

	var hostsSlice []structs.Host
	err = surrealdb.Unmarshal(hosts, &hostsSlice)

	var requestedHostByAgentID structs.Host
	for _, host := range hostsSlice {
		if host.HostID == requestedAgentByID.HostID {
			requestedHostByAgentID = host
			// FIXME: Logging
			return c.Status(fiber.StatusOK).JSON(requestedHostByAgentID)
		}
	}

	// IF HERE:
	// No Host or No agent found
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Agent Found, but no Host is associated... Thats Weird.",
	})
}
