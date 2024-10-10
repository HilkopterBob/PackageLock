package handler

import (
	"os"
	"packagelock/config"
	"packagelock/db"
	"packagelock/logger"
	"packagelock/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/surrealdb/surrealdb.go"
)

func LoginHandler(c *fiber.Ctx) error {
	// Data Sheme
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Cast POST
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		logger.Logger.Debugf("Got invalid Login Request: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	data, err := db.DB.Select("user")
	if err != nil {
		logger.Logger.Warnf("Got error while 'db.DB.Select('user')': %s \nSending 'StatusInternalServerError'", err)
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}

	var UserTable []structs.User
	err = surrealdb.Unmarshal(data, &UserTable)
	if err != nil {
		logger.Logger.Warnf("Got error while surrealdb.Unmarshal: %s \nSending 'StatusInternalServerError'", err)
		return c.Status(fiber.StatusInternalServerError).JSON(nil)
	}

	var authenticatedUser structs.User
	for _, possibleUser := range UserTable {
		// TODO: implement password hashing
		if possibleUser.Username == loginReq.Username && possibleUser.Password == loginReq.Password {
			authenticatedUser = possibleUser
		}
	}

	// create JWT
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = authenticatedUser.Username
	claims["userID"] = authenticatedUser.UserID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days expiry

	// Sign and get the encoded token
	keyData, err := os.ReadFile(config.Config.GetString("network.ssl-config.privatekeypath"))
	if err != nil {
		logger.Logger.Warnf("Can't read from Private Key File, got: %s\nPath: %s", err, config.Config.GetString("network.ssl-config.privatekeypath"))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		logger.Logger.Warnf("Can't sign with Private Key File, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		logger.Logger.Warnf("Can't generate JWT, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Create && Populate new JWT-Onject
	// for appending to User && storing in DB
	newJWT := structs.ApiKey{
		KeyValue:         tokenString,
		Description:      "User Generated JWT",
		AccessSeperation: false,
		AccessRights:     make([]string, 0),
		CreationTime:     time.Now(),
		UpdateTime:       time.Now(),
	}

	// Add the token to the user's APIToken slice
	authenticatedUser.ApiKeys = append(authenticatedUser.ApiKeys, newJWT)
	authenticatedUser.UpdateTime = time.Now()

	_, err = db.DB.Update(authenticatedUser.ID, authenticatedUser)
	if err != nil {
		logger.Logger.Warnf("Can't update User entry in DB to append new JWT, got: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	logger.Logger.Infof("Successfully authenticated and authorized following User: %s", authenticatedUser.Username)
	return c.JSON(newJWT)
}
