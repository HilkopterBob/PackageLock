package handler

import (
	"context"
	"packagelock/db"
	"packagelock/structs"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetHosts responds with a list of all hosts.
func GetHosts(c *fiber.Ctx) error {
	allHostsCursor, err := db.Client.Database("packagelock").Collection("hosts").Find(context.Background(), bson.D{})
	if err != nil {
		// TODO: Logging
		// TODO: Error handling
		return err
	}
	var result []structs.Host
	if err = allHostsCursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// GetAgents responds with a list of all agents.
func GetAgents(c *fiber.Ctx) error {
	allAgentsCursor, err := db.Client.Database("packagelock").Collection("agents").Find(context.Background(), bson.D{})
	if err != nil {
		// TODO: Logging
		// TODO: Error handling
		return err
	}
	var result []structs.Agent
	if err = allAgentsCursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}
