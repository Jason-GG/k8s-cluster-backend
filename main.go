package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/routes"
	// "github.com/sjian_mstr/cluster-management/redis_connection"
)

func main() {
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
