package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"packagelock/config"
	"packagelock/db"
	"packagelock/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
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

	// var user structs.User
	// Find the user by username
	filter := bson.D{
		{"username", loginReq.Username},
		{"password", loginReq.Password},
	}
	// usersColl := db.Client.Database("packagelock").Collection("users")
	cursor, err := db.Client.Database("packagelock").Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	var result []structs.User
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	fmt.Println(filter)
	fmt.Println(result)
	// As 'user' is a struct, check for a must-have value (USerID)
	// If UserID == "" the user couldn't be found -> doesn't exist!
	//if user.UserID == "" {
	//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//		"error": "Invalid username or password",
	//	})
	//}

	// Validate the password
	// Generate JWT token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = result[0].Username
	claims["userID"] = result[0].UserID
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
		log.Println("Failed to generate JWT token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	newJWT := structs.ApiKey{
		KeyValue:         tokenString,
		Description:      "User Generated JWT",
		AccessSeperation: false,
		AccessRights:     make([]string, 0),
		CreationTime:     time.Now(),
		UpdateTime:       time.Now(),
	}

	// Add the token to the user's APIToken slice
	result[0].ApiKeys = append(result[0].ApiKeys, newJWT)
	filter = bson.D{
		{"userid", result[0].UserID},
	}

	result[0].UpdateTime = time.Now()

	replacement, err := bson.Marshal(result[0])
	if err != nil {
		return err
	}
	fmt.Println("HALLO")
	updateResult, err := db.Client.Database("packagelock").Collection("users").ReplaceOne(context.Background(), filter, replacement)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(updateResult)

	fmt.Println("Printing Result now:")
	fmt.Println(result[0])

	return c.JSON(newJWT)
}
