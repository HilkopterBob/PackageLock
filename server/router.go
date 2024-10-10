package server

import (
	"log"
	"os"
	"packagelock/config"
	"packagelock/handler"
	"packagelock/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

// Routes holds the Fiber app instance.
type Routes struct {
	Router *fiber.App
}

// addAgentHandler sets up agent-related routes in Fiber.
func (r Routes) addAgentHandler(group fiber.Router) {
	AgentGroup := group.Group("/agents")

	AgentGroup.Get("/", handler.GetAgentByID)
	AgentGroup.Post("/register", handler.RegisterAgent)
	logger.Logger.Debug("Added Agent Handlers.")
}

// addGeneralHandler sets up general-related routes in Fiber.
func (r Routes) addGeneralHandler(group fiber.Router) {
	GeneralGroup := group.Group("/general")

	GeneralGroup.Get("/hosts", handler.GetHosts)
	GeneralGroup.Get("/agents", handler.GetAgents)
	logger.Logger.Debug("Added General Handlers.")
}

// addHostHandler sets up host-related routes in Fiber.
func (r Routes) addHostHandler(group fiber.Router) {
	HostGroup := group.Group("/hosts")

	HostGroup.Get("/", handler.GetHostByAgentID)
	HostGroup.Post("/register", handler.RegisterHost)

	logger.Logger.Debug("Added Host Handlers.")
}

func (r Routes) addLoginHandler(group fiber.Router) {
	LoginGroup := group.Group("/auth")

	LoginGroup.Post("/login", handler.LoginHandler)

	logger.Logger.Debug("Added Login Handlers.")
}

// AddRoutes adds all handler groups to the current Fiber app.
// It's exported and used in main() to return the configured Router.
func AddRoutes(Config config.ConfigProvider) Routes {
	// Initialize template engine
	engine := html.New("./templates", ".html")

	// Initialize Fiber app
	router := Routes{
		Router: fiber.New(fiber.Config{
			Views: engine,
		}),
	}

	router.addLoginHandler(router.Router)

	// Use JWT if in production
	if Config.Get("general.production") == true {
		logger.Logger.Info("Enabled Production! Adding JWT!")
		// Read the private key for JWT
		keyData, err := os.ReadFile(Config.GetString("network.ssl.privatekeypath"))
		if err != nil {
			log.Fatal(err)
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			logger.Logger.Panicf("Can't open Private Key File! Got: %s", err)
		}

		// JWT Middleware to protect specific routes
		jwtMiddleware := jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: privateKey},
		})

		// Apply JWT protection to all routes in the "/v1" group
		v1 := router.Router.Group("/v1", jwtMiddleware)

		// Add route handlers to the protected group
		router.addGeneralHandler(v1)
		router.addAgentHandler(v1)
		router.addHostHandler(v1)
	} else {
		logger.Logger.Info("Non-Production Setup! Disabled JWT!")
		// Create the versioned route group without JWT protection (for non-production environments)
		v1 := router.Router.Group("/v1")

		// Add route handlers without JWT protection
		router.addGeneralHandler(v1)
		router.addAgentHandler(v1)
		router.addHostHandler(v1)
	}

	// Middleware to recover from panics
	router.Router.Use(recover.New())

	// Add 404 handler
	router.Router.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{})
	})

	return router
}

// ListenAndServeTLS starts the Fiber server using TLS (HTTPS)
func ListenAndServeTLS(router *fiber.App, certFile, keyFile, addr string) error {
	// Start HTTPS server using the provided certificate and key files
	return router.ListenTLS(addr, certFile, keyFile)
}
