package routes

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/controller"
)

func getPodLogs(c *fiber.Ctx) error {
	queryParams := c.Queries()
	env := queryParams["env"]
	namespace := queryParams["namespace"]
	podName := queryParams["podName"]
	number := queryParams["number"]
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	podLogs, err := controller.GetPodLogs(env, namespace, podName, num)
	if err != nil {
		fmt.Printf("Error pod logs: %v\n", err)
		// os.Exit(1)
		return jsonResponse(c, fiber.StatusOK, "OK", 1, err)
	}

	return jsonResponse(c, fiber.StatusOK, "OK", 1, podLogs)
}
