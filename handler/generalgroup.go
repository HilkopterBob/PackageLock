package handler

import (
	"github.com/gofiber/fiber/v2"
)

// GetHosts responds with a list of all hosts.
func GetHosts(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Hosts)
}

// GetAgents responds with a list of all agents.
func GetAgents(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Agents)
}
