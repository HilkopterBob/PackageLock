package handler

import (
	"encoding/base64"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

// GetAgentByID filters a slice of Agents for a matching Agent.Agent_ID.
// It returns a JSON response with fiber.StatusOK or fiber.StatusNotFound.
func GetAgentByID(c *fiber.Ctx) error {
	// ID is an URL slice. Its a URL-Save base64 encoded UUID
	urlIDBytes, err := base64.RawURLEncoding.DecodeString(c.Query("AgentID"))
	if err != nil {
		logger.Logger.Warnf("Can't parse AgentID from URL, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse AgentID.",
		})
	}

	urlIDString := string(urlIDBytes)

	agents, err := db.DB.Select("agents")
	if err != nil {
		logger.Logger.Warnf("Failed to fetch 'agents' from db, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch Agents.",
		})
	}

	var agentsSlice []structs.Agent
	err = surrealdb.Unmarshal(agents, &agentsSlice)
	if err != nil {
		logger.Logger.Warnf("Failed to unmarshal agents, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed Unmarshal.",
		})
	}

	var requestedAgentByID structs.Agent
	for _, agent := range agentsSlice {
		if agent.AgentID.String() == urlIDString {
			requestedAgentByID = agent
			return c.Status(fiber.StatusOK).JSON(requestedAgentByID)
		}
	}

	logger.Logger.Warnf("Got Request for agent with id: %s, which dosn't exist!", urlIDString)
	return c.Status(fiber.StatusNotFound).JSON(nil)
}

// RegisterAgent handles POST requests to register a new agent.
func RegisterAgent(c *fiber.Ctx) error {
	var newAgent structs.Agent

	// Parse the JSON request body into newAgent
	if err := c.BodyParser(&newAgent); err != nil {
		logger.Logger.Warnf("Cannot parse JSON into new Agent! Got: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	newAgentInsertionData, err := db.DB.Create("agents", newAgent)
	if err != nil {
		logger.Logger.Warnf("Can't insert new Agent into DB, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}
	// Respond with the newly created agent
	logger.Logger.Infof("Successfully Created new Agent with ID: %s", newAgent.AgentID)
	return c.Status(fiber.StatusCreated).JSON(newAgentInsertionData)
}

// GetHostByAgentID finds the host for a given agent ID.
func GetHostByAgentID(c *fiber.Ctx) error {
	urlIDBytes, err := base64.RawURLEncoding.DecodeString(c.Query("AgentID"))
	if err != nil {
		logger.Logger.Warnf("Can't parse AgentID from URL, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse AgentID.",
		})
	}

	urlIDString := string(urlIDBytes)

	agents, err := db.DB.Select("agents")
	if err != nil {
		logger.Logger.Warnf("Failed to fetch 'agents' from db, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch Agents.",
		})
	}

	var agentsSlice []structs.Agent
	err = surrealdb.Unmarshal(agents, &agentsSlice)
	if err != nil {
		logger.Logger.Warnf("Failed to unmarshal agents, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed Unmarshal.",
		})
	}

	var requestedAgentByID structs.Agent
	for _, agent := range agentsSlice {
		if agent.AgentID.String() == urlIDString {
			requestedAgentByID = agent
		}
	}

	// if no matching agent is found, the ID is empty
	if requestedAgentByID.ID == "" {
		logger.Logger.Warnf("Got Request for agent with id: %s, which dosn't exist!", urlIDString)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Agent not found",
		})
	}

	hosts, err := db.DB.Select("hosts")
	if err != nil {
		logger.Logger.Warnf("Failed to fetch 'hosts' from db, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch Hosts.",
		})
	}

	var hostsSlice []structs.Host
	err = surrealdb.Unmarshal(hosts, &hostsSlice)
	if err != nil {
		logger.Logger.Warnf("Failed to unmarshal hosts, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed Unmarshal.",
		})
	}

	var requestedHostByAgentID structs.Host
	for _, host := range hostsSlice {
		if host.HostID == requestedAgentByID.HostID {
			requestedHostByAgentID = host
			return c.Status(fiber.StatusOK).JSON(requestedHostByAgentID)
		}
	}

	logger.Logger.Warnf("Agent Found, but no Host is associated... Thats Weird. AgentID: %s", requestedAgentByID.AgentID)
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Agent Found, but no Host is associated... Thats Weird.",
	})
}
