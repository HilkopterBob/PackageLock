package server

import (
	"packagelock/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

// Routes holds the Fiber app instance.
type Routes struct {
	Router *fiber.App
}

// addAgentHandler sets up agent-related routes in Fiber.
func (r Routes) addAgentHandler(group fiber.Router) {
	AgentGroup := group.Group("/agent")

	AgentGroup.Get("/:id", handler.GetAgentByID)
	AgentGroup.Get("/:id/host", handler.GetHostByAgentID)
	AgentGroup.Post("/register", handler.RegisterAgent)
}

// addGeneralHandler sets up general-related routes in Fiber.
func (r Routes) addGeneralHandler(group fiber.Router) {
	GeneralGroup := group.Group("/general")

	GeneralGroup.Get("/hosts", handler.GetHosts)
	GeneralGroup.Get("/agents", handler.GetAgents)
}

// addHostHandler sets up host-related routes in Fiber.
func (r Routes) addHostHandler(group fiber.Router) {
	HostGroup := group.Group("/host")

	HostGroup.Post("/register", handler.RegisterHost)
}

// AddRoutes adds all handler groups to the current Fiber app.
// It's exported and used in main() to return the configured Router.
func AddRoutes() Routes {
	// Initialize template engine
	engine := html.New("./templates", ".html")

	// Initialize Fiber app
	router := Routes{
		Router: fiber.New(fiber.Config{
			Views: engine,
		}),
	}

	router.Router.Use(recover.New())

	// Create the versioned route group
	v1 := router.Router.Group("/v1")

	// Add all route handlers
	router.addGeneralHandler(v1)
	router.addAgentHandler(v1)
	router.addHostHandler(v1)

	router.Router.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{})
	})

	return router
}
