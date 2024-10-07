package handler

import (
	"log"
	"os"
	"packagelock/config"
	"packagelock/db"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	data, err := db.DB.Select("user")
	if err != nil {
		// FIXME: error handling
		panic(err)
	}

	var UserTable []structs.User
	err = surrealdb.Unmarshal(data, &UserTable)
	if err != nil {
		// FIXME: error handling
		panic(err)
	}

	var authenticatedUser structs.User
	for _, possibleUser := range UserTable {
		// TODO: implement password hashing
		if possibleUser.Username == loginReq.Username && possibleUser.Password == loginReq.Password {
			authenticatedUser = possibleUser
			// TODO: log token creation
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
		log.Fatal(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Fatal(err)
	}
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		// FIXME: error handling & maybe logging to?
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
		// FIXME: errorhandling
		panic(err)
	}

	return c.JSON(newJWT)
}
