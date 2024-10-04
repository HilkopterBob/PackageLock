package handler

import (
	"packagelock/db"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
	"github.com/surrealdb/surrealdb.go"
)

// GetHosts responds with a list of all hosts.
func GetHosts(c *fiber.Ctx) error {
	hosts, err := db.DB.Select("hosts")
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}
	var hostsSlice []structs.Host
	err = surrealdb.Unmarshal(hosts, &hostsSlice)
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}
	return c.Status(fiber.StatusOK).JSON(hostsSlice)
}

// GetAgents responds with a list of all agents.
func GetAgents(c *fiber.Ctx) error {
	agents, err := db.DB.Select("agents")
	if err != nil {
		// FIXME: Error handling
		panic(err)
	}
	var agentsSlice []structs.Host
	err = surrealdb.Unmarshal(agents, &agentsSlice)
	if err != nil {
		// FIXME: error handling
		panic(err)
	}
	return c.Status(fiber.StatusOK).JSON(agentsSlice)
}
