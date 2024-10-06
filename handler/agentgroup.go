package handler

import (
	"context"
	"fmt"
	"packagelock/db"
	"packagelock/structs"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAgentByID filters a slice of Agents for a matching Agent.Agent_ID.
// It returns a JSON response with fiber.StatusOK or fiber.StatusNotFound.
func GetAgentByID(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		fmt.Println(err)
	}

	filter := bson.D{
		{"agent_id", id},
	}

	fmt.Println(filter)

	AgentsCursor, err := db.Client.Database("packagelock").Collection("agents").Find(context.Background(), filter)
	if err != nil {
		// TODO: Logging
		// TODO: Error handling
		return err
	}

	fmt.Println(AgentsCursor)

	var agents []structs.Agent
	if err = AgentsCursor.All(context.Background(), &agents); err != nil {
		panic(err)
	}

	fmt.Println(agents)

	if len(agents) <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no agent under that id"})
	}

	return c.Status(fiber.StatusOK).JSON(agents)
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

	coll := db.Client.Database("packagelock").Collection("agents")
	_, err := coll.InsertOne(context.Background(), newAgent)
	if err != nil {
		return fmt.Errorf("failed to add new Agent to db: %w", err)
	}

	// Respond with the newly created agent
	return c.Status(fiber.StatusCreated).JSON(newAgent)
}

// GetHostByAgentID finds the host for a given agent ID.
func GetHostByAgentID(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		fmt.Println(err)
	}
	filter := bson.D{
		{"agent_id", id},
	}

	AgentsCursor, err := db.Client.Database("packagelock").Collection("agents").Find(context.Background(), filter)
	if err != nil {
		// TODO: Logging
		// TODO: Error handling
		return err
	}

	var agents []structs.Agent
	if err = AgentsCursor.All(context.TODO(), &agents); err != nil {
		panic(err)
	}

	if len(agents) <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no agent with that id"})
	}

	// TODO: Filter all hosts for agent id
	filter = bson.D{
		{"id", agents[0].HostID},
	}

	fmt.Println(filter)

	HostsCursor, err := db.Client.Database("packagelock").Collection("hosts").Find(context.Background(), filter)
	var hosts []structs.Agent
	if err = HostsCursor.All(context.TODO(), &hosts); err != nil {
		panic(err)
	}

	if len(hosts) <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "no host with that id", "agent": agents[0]})
	}

	return c.Status(fiber.StatusOK).JSON(hosts)
}
