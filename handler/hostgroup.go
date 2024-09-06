package handler

import (
	"context"
	"fmt"
	"packagelock/db"
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

	coll := db.Client.Database("packagelock").Collection("hosts")
	_, err := coll.InsertOne(context.Background(), newHost)
	if err != nil {
		return fmt.Errorf("failed to add new Host to db: %w", err)
	}

	// Respond with the newly created agent
	return c.Status(fiber.StatusCreated).JSON(newHost)
}
