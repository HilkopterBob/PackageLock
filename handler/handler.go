package handler

import (
	"encoding/base64"
	"os"
	"packagelock/db"
	"packagelock/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/surrealdb/surrealdb.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handlers struct {
	// UserGroup handlers
	LoginHandler fiber.Handler

	// AgentGroup handlers
	GetAgentByID     fiber.Handler
	RegisterAgent    fiber.Handler
	GetHostByAgentID fiber.Handler

	// GeneralGroup handlers
	GetHosts  fiber.Handler
	GetAgents fiber.Handler

	// HostGroup handlers
	RegisterHost fiber.Handler
}

type HandlerParams struct {
	fx.In

	Logger *zap.Logger
	Config *viper.Viper
	DB     *db.Database
}

// NewHandlers constructs all handler functions with injected dependencies.
func NewHandlers(params HandlerParams) *Handlers {
	return &Handlers{
		LoginHandler:     NewLoginHandler(params),
		GetAgentByID:     NewGetAgentByIDHandler(params),
		RegisterAgent:    NewRegisterAgentHandler(params),
		GetHostByAgentID: NewGetHostByAgentIDHandler(params),
		GetHosts:         NewGetHostsHandler(params),
		GetAgents:        NewGetAgentsHandler(params),
		RegisterHost:     NewRegisterHostHandler(params),
	}
}

func NewLoginHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var loginReq LoginRequest
		if err := c.BodyParser(&loginReq); err != nil {
			params.Logger.Debug("Invalid login request", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse request",
			})
		}

		data, err := params.DB.DB.Select("user")
		if err != nil {
			params.Logger.Warn("Error selecting 'user'", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(nil)
		}

		var userTable []structs.User
		err = surrealdb.Unmarshal(data, &userTable)
		if err != nil {
			params.Logger.Warn("Error unmarshalling users", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(nil)
		}

		var authenticatedUser *structs.User
		for _, possibleUser := range userTable {
			// TODO: Implement password hashing
			if possibleUser.Username == loginReq.Username && possibleUser.Password == loginReq.Password {
				authenticatedUser = &possibleUser
				break
			}
		}

		if authenticatedUser == nil {
			// User not found or password incorrect
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}

		// Create JWT
		token := jwt.New(jwt.SigningMethodRS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = authenticatedUser.Username
		claims["userID"] = authenticatedUser.UserID
		claims["exp"] = time.Now().Add(72 * time.Hour).Unix() // 3 days expiry

		// Sign and get the encoded token
		keyData, err := os.ReadFile(params.Config.GetString("network.ssl-config.privatekeypath"))
		if err != nil {
			params.Logger.Warn("Cannot read private key file", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			params.Logger.Warn("Cannot parse private key", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		tokenString, err := token.SignedString(privateKey)
		if err != nil {
			params.Logger.Warn("Cannot generate JWT", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		// Create and append new JWT object
		newJWT := structs.ApiKey{
			KeyValue:         tokenString,
			Description:      "User Generated JWT",
			AccessSeperation: false,
			AccessRights:     []string{},
			CreationTime:     time.Now(),
			UpdateTime:       time.Now(),
		}

		// Update user with new API key
		authenticatedUser.ApiKeys = append(authenticatedUser.ApiKeys, newJWT)
		authenticatedUser.UpdateTime = time.Now()

		_, err = params.DB.DB.Update(authenticatedUser.ID, authenticatedUser)
		if err != nil {
			params.Logger.Warn("Cannot update user in DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		params.Logger.Info("User authenticated", zap.String("username", authenticatedUser.Username))
		return c.JSON(newJWT)
	}
}

func NewGetAgentByIDHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		urlIDBytes, err := base64.RawURLEncoding.DecodeString(c.Query("AgentID"))
		if err != nil {
			params.Logger.Warn("Cannot parse AgentID from URL", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse AgentID",
			})
		}

		urlIDString := string(urlIDBytes)

		agents, err := params.DB.DB.Select("agents")
		if err != nil {
			params.Logger.Warn("Failed to fetch 'agents' from DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch agents",
			})
		}

		var agentsSlice []structs.Agent
		err = surrealdb.Unmarshal(agents, &agentsSlice)
		if err != nil {
			params.Logger.Warn("Failed to unmarshal agents", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to unmarshal agents",
			})
		}

		var requestedAgent *structs.Agent
		for _, agent := range agentsSlice {
			if agent.AgentID.String() == urlIDString {
				requestedAgent = &agent
				break
			}
		}

		if requestedAgent == nil {
			params.Logger.Warn("Agent not found", zap.String("AgentID", urlIDString))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Agent not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(requestedAgent)
	}
}

func NewRegisterAgentHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newAgent structs.Agent

		if err := c.BodyParser(&newAgent); err != nil {
			params.Logger.Warn("Cannot parse JSON into new Agent", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		newAgentInsertionData, err := params.DB.DB.Create("agents", newAgent)
		if err != nil {
			params.Logger.Warn("Cannot insert new Agent into DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(nil)
		}

		params.Logger.Info("Created new Agent", zap.String("AgentID", newAgent.AgentID.String()))
		return c.Status(fiber.StatusCreated).JSON(newAgentInsertionData)
	}
}

func NewGetHostByAgentIDHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		urlIDBytes, err := base64.RawURLEncoding.DecodeString(c.Query("AgentID"))
		if err != nil {
			params.Logger.Warn("Cannot parse AgentID from URL", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse AgentID",
			})
		}

		urlIDString := string(urlIDBytes)

		agents, err := params.DB.DB.Select("agents")
		if err != nil {
			params.Logger.Warn("Failed to fetch 'agents' from DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch agents",
			})
		}

		var agentsSlice []structs.Agent
		err = surrealdb.Unmarshal(agents, &agentsSlice)
		if err != nil {
			params.Logger.Warn("Failed to unmarshal agents", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to unmarshal agents",
			})
		}

		var requestedAgent *structs.Agent
		for _, agent := range agentsSlice {
			if agent.AgentID.String() == urlIDString {
				requestedAgent = &agent
				break
			}
		}

		if requestedAgent == nil {
			params.Logger.Warn("Agent not found", zap.String("AgentID", urlIDString))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Agent not found",
			})
		}

		hosts, err := params.DB.DB.Select("hosts")
		if err != nil {
			params.Logger.Warn("Failed to fetch 'hosts' from DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch hosts",
			})
		}

		var hostsSlice []structs.Host
		err = surrealdb.Unmarshal(hosts, &hostsSlice)
		if err != nil {
			params.Logger.Warn("Failed to unmarshal hosts", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to unmarshal hosts",
			})
		}

		var requestedHost *structs.Host
		for _, host := range hostsSlice {
			if host.HostID == requestedAgent.HostID {
				requestedHost = &host
				break
			}
		}

		if requestedHost == nil {
			params.Logger.Warn("No host associated with agent", zap.String("AgentID", requestedAgent.AgentID.String()))
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Host not found for the agent",
			})
		}

		return c.Status(fiber.StatusOK).JSON(requestedHost)
	}
}

func NewGetHostsHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hosts, err := params.DB.DB.Select("hosts")
		if err != nil {
			params.Logger.Warn("Failed to fetch 'hosts' from DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch hosts",
			})
		}

		var hostsSlice []structs.Host
		err = surrealdb.Unmarshal(hosts, &hostsSlice)
		if err != nil {
			params.Logger.Warn("Failed to unmarshal hosts", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to unmarshal hosts",
			})
		}

		return c.Status(fiber.StatusOK).JSON(hostsSlice)
	}
}

func NewGetAgentsHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		agents, err := params.DB.DB.Select("agents")
		if err != nil {
			params.Logger.Warn("Failed to fetch 'agents' from DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch agents",
			})
		}

		var agentsSlice []structs.Agent
		err = surrealdb.Unmarshal(agents, &agentsSlice)
		if err != nil {
			params.Logger.Warn("Failed to unmarshal agents", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to unmarshal agents",
			})
		}

		return c.Status(fiber.StatusOK).JSON(agentsSlice)
	}
}

func NewRegisterHostHandler(params HandlerParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newHost structs.Host

		if err := c.BodyParser(&newHost); err != nil {
			params.Logger.Warn("Cannot parse JSON into new Host", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		transaction, err := params.DB.DB.Create("hosts", newHost)
		if err != nil {
			params.Logger.Warn("Cannot insert new Host into DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(nil)
		}

		params.Logger.Info("Created new Host", zap.String("HostID", newHost.HostID.String()))
		return c.Status(fiber.StatusCreated).JSON(transaction)
	}
}

// Module exports the handlers as an Fx module.
var Module = fx.Options(
	fx.Provide(NewHandlers),
)
