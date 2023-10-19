package database_table

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/database"
	"github.com/sjian_mstr/cluster-management/models"
)

var (
	jwtSecretKey  = []byte("your-secret-key") // Change this to your secret key
	tokenDuration = time.Hour * 24            // Token expiration duration
)

func LoginHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	if err := database.Database.Db.Where("username = ?", requestBody.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Decode the stored password from Base64
	storedPasswordBytes, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode password"})
	}
	storedPassword := string(storedPasswordBytes)

	// Check if the provided password matches the user's password
	if requestBody.Password != storedPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate a JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(tokenDuration).Unix()

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Set the token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(tokenDuration),
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Return the token to the client
	response := struct {
		Data string `json:"data"`
		Msg  string `json:"msg"`
		Code int    `json:"code"`
	}{
		Data: "Login successful",
		Msg:  "ok",
		Code: 1,
	}
	return c.JSON(response)
}

func AuthMiddleware(c *fiber.Ctx) error {
	// Retrieve the "Authorization" header containing basic auth credentials
	authHeader := c.Get("Authorization")

	// Check if the header is present and contains the "Basic " prefix
	if authHeader != "" && strings.HasPrefix(authHeader, "Basic ") {
		// Extract the base64-encoded credentials
		credentialsBase64 := strings.TrimPrefix(authHeader, "Basic ")
		credentialsBytes, err := base64.StdEncoding.DecodeString(credentialsBase64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Split the credentials into username and password
		credentials := strings.SplitN(string(credentialsBytes), ":", 2)
		if len(credentials) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		username := credentials[0]
		password := credentials[1]

		// Look up the user in the database
		var user models.User
		if err := database.Database.Db.Where("username = ?", username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Decode the stored password from Base64
		storedPasswordBytes, err := base64.StdEncoding.DecodeString(user.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode password"})
		}
		storedPassword := string(storedPasswordBytes)

		// Check if the provided password matches the user's password
		if password != storedPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Next()
	}

	// If Basic Auth is not provided, check for the "auth_token" cookie
	tokenString := c.Cookies("auth_token")

	// Check if the token is missing or improperly formatted
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Continue to the next middleware or route
	return c.Next()
}
