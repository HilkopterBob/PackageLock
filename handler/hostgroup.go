package handler

import (
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
)

// RegisterHost handles the registration of a new host.
func RegisterHost(c *fiber.Ctx) error {
	var newHost structs.Host

	// Parse the JSON request body into newHost
	if err := c.BodyParser(&newHost); err != nil {
		// TODO: Add logs
		// TODO: Add error handling
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Append new host to the Hosts slice
	Hosts = append(Hosts, newHost)

	// Respond with the newly created host
	return c.Status(fiber.StatusCreated).JSON(newHost)
}
