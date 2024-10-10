package handler

import (
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
)

// RegisterHost handles the registration of a new host.
func RegisterHost(c *fiber.Ctx) error {
	var newHost structs.Host

	// Parse the JSON request body into newHost
	if err := c.BodyParser(&newHost); err != nil {
		logger.Logger.Warnf("Cannot parse JSON into new Host! Got: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	transaction, err := db.DB.Create("hosts", newHost)
	if err != nil {
		logger.Logger.Warnf("Can't insert new Host into DB, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}

	logger.Logger.Infof("Successfully Created new Host with ID: %s", newHost.HostID)
	return c.Status(fiber.StatusCreated).JSON(transaction)
}
