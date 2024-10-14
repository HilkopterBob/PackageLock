package server

import (
	"context"
	"os"
	"packagelock/handler"
	"strconv"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Config    *viper.Viper
	Handlers  *handler.Handlers // The injected Handlers struct
}

func NewServer(params ServerParams, logger *zap.Logger) *fiber.App {
	logger.Info("Starting API-Server Initialization:")
	// Initialize template engine
	engine := html.New("./templates", ".html")
	logger.Info("Added template Enginge.")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware to recover from panics
	app.Use(recover.New())
	logger.Info("Added Recovery Middleware.")

	// Add routes
	addRoutes(app, params)
	logger.Info("Added routes.")

	appVersion := params.Config.GetString("general.app-version")

	// Add 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{
			"AppVersion": appVersion,
		})
	})
	logger.Info("Added default 404 Handler.")

	// Start the server using lifecycle hooks
	logger.Info("Finished API-Server Initialization.")
	startServer(app, params)

	return app
}

func addRoutes(app *fiber.App, params ServerParams) {
	// Add login handler
	addLoginHandler(app, params)

	// Use JWT if in production
	if params.Config.GetBool("general.production") {
		params.Logger.Info("Enabled Production! Adding JWT!")

		// Read the private key for JWT
		keyData, err := os.ReadFile(params.Config.GetString("network.ssl-config.privatekeypath"))
		if err != nil {
			params.Logger.Fatal("Failed to read private key for JWT", zap.Error(err))
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			params.Logger.Fatal("Failed to parse private key for JWT", zap.Error(err))
		}

		// JWT Middleware to protect specific routes
		jwtMiddleware := jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: privateKey},
		})

		// Apply JWT protection to all routes in the "/v1" group
		v1 := app.Group("/v1", jwtMiddleware)

		// Add route handlers to the protected group
		addGeneralHandler(v1, params)
		addAgentHandler(v1, params)
		addHostHandler(v1, params)
	} else {
		params.Logger.Info("Non-Production Setup! Disabled JWT!")

		// Create the versioned route group without JWT protection
		v1 := app.Group("/v1")

		// Add route handlers without JWT protection
		addGeneralHandler(v1, params)
		addAgentHandler(v1, params)
		addHostHandler(v1, params)
	}
}

func startServer(app *fiber.App, params ServerParams) {
	serverAddr := params.Config.GetString("network.fqdn") + ":" + params.Config.GetString("network.port")

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Write PID to file
			pid := os.Getpid()
			err := os.WriteFile("packagelock.pid", []byte(strconv.Itoa(pid)), 0644)
			if err != nil {
				params.Logger.Warn("Failed to write PID file", zap.Error(err))
			} else {
				params.Logger.Info("PID file written", zap.Int("PID", pid))
			}

			go func() {
				if params.Config.GetBool("network.ssl") {
					params.Logger.Info("Starting HTTPS server", zap.String("address", serverAddr))

					certFile := params.Config.GetString("network.ssl-config.certificatepath")
					keyFile := params.Config.GetString("network.ssl-config.privatekeypath")

					if err := app.ListenTLS(serverAddr, certFile, keyFile); err != nil {
						params.Logger.Fatal("Failed to start HTTPS server", zap.Error(err))
					}
				} else {
					params.Logger.Info("Starting HTTP server", zap.String("address", serverAddr))

					if err := app.Listen(serverAddr); err != nil {
						params.Logger.Fatal("Failed to start HTTP server", zap.Error(err))
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Shutting down server")
			return app.Shutdown()
		},
	})
}

// Individual route handler functions
func addAgentHandler(group fiber.Router, params ServerParams) {
	agentGroup := group.Group("/agents")

	agentGroup.Get("/", params.Handlers.GetAgentByID)
	agentGroup.Post("/register", params.Handlers.RegisterAgent)
	params.Logger.Debug("Added Agent Handlers.")
}

func addGeneralHandler(group fiber.Router, params ServerParams) {
	generalGroup := group.Group("/general")

	generalGroup.Get("/hosts", params.Handlers.GetHosts)
	generalGroup.Get("/agents", params.Handlers.GetAgents)
	params.Logger.Debug("Added General Handlers.")
}

func addHostHandler(group fiber.Router, params ServerParams) {
	hostGroup := group.Group("/hosts")

	hostGroup.Get("/", params.Handlers.GetHostByAgentID)
	hostGroup.Post("/register", params.Handlers.RegisterHost)
	params.Logger.Debug("Added Host Handlers.")
}

func addLoginHandler(group fiber.Router, params ServerParams) {
	loginGroup := group.Group("/auth")

	loginGroup.Post("/login", params.Handlers.LoginHandler)
	params.Logger.Debug("Added Login Handlers.")
}

// Module exports the server module.
var Module = fx.Options(
	fx.Provide(NewServer),
)
