package handler

import (
	"log"
	"os"
	"packagelock/config"
	"packagelock/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	var user structs.User
	// Find the user by username
	for _, u := range Users {
		if u.Username == loginReq.Username {
			user = u
			break
		}
	}

	// As 'user' is a struct, check for a must-have value (USerID)
	// If UserID == "" the user couldn't be found -> doesn't exist!
	if user.UserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Validate the password
	if user.Password != loginReq.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["userID"] = user.UserID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days expiry

	// Sign and get the encoded token
	keyData, err := os.ReadFile(config.Config.GetString("network.ssl.privatekeypath"))
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Fatal(err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Println("Failed to generate JWT token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Add the token to the user's APIToken slice
	user.APIToken = append(user.APIToken, tokenString)

	// Return the token and user information
	return c.JSON(fiber.Map{
		"message":  "Login successful",
		"token":    tokenString,
		"username": user.Username,
	})
}
