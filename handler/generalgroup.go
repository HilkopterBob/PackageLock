package handler

import (
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

// GetHosts responds with a list of all hosts.
func GetHosts(c *fiber.Ctx) error {
	hosts, err := db.DB.Select("hosts")
	if err != nil {
		logger.Logger.Warnf("Failed to fetch 'hosts' from db, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch hosts.",
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

	return c.Status(fiber.StatusOK).JSON(hostsSlice)
}

// GetAgents responds with a list of all agents.
func GetAgents(c *fiber.Ctx) error {
	agents, err := db.DB.Select("agents")
	if err != nil {
		logger.Logger.Warnf("Failed to fetch 'agents' from db, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch agents.",
		})
	}

	var agentsSlice []structs.Host
	err = surrealdb.Unmarshal(agents, &agentsSlice)
	if err != nil {
		logger.Logger.Warnf("Failed to fetch 'agents' from db, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch agents.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(agentsSlice)
}
