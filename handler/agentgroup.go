package handler

import (
	"packagelock/structs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetAgentByID filters a slice of Agents for a matching Agent.Agent_ID.
// It returns a JSON response with fiber.StatusOK or fiber.StatusNotFound.
func GetAgentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	for _, a := range Agents {
		if strconv.Itoa(a.Host_ID) == id {
			return c.Status(fiber.StatusOK).JSON(a)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no agent under that id"})
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

	// Add new agent to the Agents slice
	Agents = append(Agents, newAgent)

	// Respond with the newly created agent
	return c.Status(fiber.StatusCreated).JSON(newAgent)
}

// GetHostByAgentID finds the host for a given agent ID.
func GetHostByAgentID(c *fiber.Ctx) error {
	var agentByID structs.Agent

	// Get the value from /agent/:id/host
	id := c.Params("id")

	// Find the agent by the URL-ID
	for _, a := range Agents {
		if strconv.Itoa(a.Host_ID) == id {
			agentByID = a
			break
		}
	}

	if agentByID.Agent_ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no agent under that id"})
	}

	// Find the host with the same id as the agent
	for _, host := range Hosts {
		if host.ID == agentByID.Agent_ID {
			return c.Status(fiber.StatusOK).JSON(host)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no host found for that agent"})
}
