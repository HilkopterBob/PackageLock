package handler

import (
	"packagelock/db"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
)

// RegisterHost handles the registration of a new host.
func RegisterHost(c *fiber.Ctx) error {
	var newHost structs.Host

	// Parse the JSON request body into newHost
	if err := c.BodyParser(&newHost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	transaction, err := db.DB.Create("hosts", newHost)
	if err != nil {
		// FIXME: error handling
		panic(err)
	}

	// FIXME: Logging!
	return c.Status(fiber.StatusCreated).JSON(transaction)
}
