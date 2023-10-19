package routes

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sjian_mstr/cluster-management/controller"
	"github.com/sjian_mstr/cluster-management/database_table"
)

type apiResponseFormat struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func SetupRoutes(app *fiber.App) {
	// CROS DOMAIN CONFIG
	corsConfig := cors.Config{
		AllowOrigins:     "https://cmdbservice.cloud.microstrategy.com", // Allow any origin
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, X-Requested-With, Content-Type, Accept, authorization",
		AllowCredentials: true,
	}
	app.Use(cors.New(corsConfig))
	app.Post("/login", database_table.LoginHandler)
	// Define a route for the health check
	app.Get("/api/healthCheck", func(c *fiber.Ctx) error {
		// You can perform any health checks here
		// For example, check if your database is reachable, external services are available, etc.
		// If everything is healthy, respond with a 200 OK status
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})
	app.Get("/ws", websocket.New(controller.ExecuteCommand))

	app.Use(database_table.AuthMiddleware)
	app.Get("/namespace", getNamespaceListInfo)
	app.Get("/varifyNamespace", verifyNamespaceInfo)
	app.Get("/listNamespace", getNamespaceListDetailInfo)
	app.Get("/deploymentList", getDeploymentListInfo)
	app.Get("/listPods", getPodsListInfo)
	app.Get("/podsListFromDeployment", getDeploymentPods)
	app.Get("/podsLogs", getPodLogs)
	app.Get("/podDelete", deletePods)
	app.Get("/listJobs", listJods)

}

func jsonResponse(c *fiber.Ctx, status int, msg string, code int, data interface{}) error {
	response := apiResponseFormat{
		Msg:  msg,
		Code: code,
		Data: data,
	}

	return c.Status(status).JSON(response)
}
