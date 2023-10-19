package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/controller"
)

func listJods(c *fiber.Ctx) error {
	queryParams := c.Queries()
	// Retrieve specific query parameters
	env := queryParams["env"]
	namespace := queryParams["namespace"]

	jodsInfo, err := controller.ListJodsInfo(env, namespace)
	if err != nil {
		fmt.Printf("Error  GetDeploymentPods: %v\n", err)
		// os.Exit(1)
		return jsonResponse(c, fiber.StatusOK, "Error", 1, err)
	}
	for _, pods := range jodsInfo.Items {
		fmt.Printf("Name: %s\n", pods.ObjectMeta.Name)
		fmt.Println("-----------------------")
	}
	// return c.Status(200).JSON(podsInfo)
	return jsonResponse(c, fiber.StatusOK, "OK", 1, jodsInfo)
}
