package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/controller"
)

func getPodsListInfo(c *fiber.Ctx) error {
	queryParams := c.Queries()
	env := queryParams["env"]
	namespace := queryParams["namespace"]

	podsInfo, err := controller.ListPodsInfo(env, namespace)
	if err != nil {
		fmt.Printf("Error listing podsInfo: %v\n", err)
		// os.Exit(1)
		return c.Status(200).JSON("Error listing podsInfo")
	}
	for _, deployment := range podsInfo {
		fmt.Printf("Name: %s\n", deployment.ObjectMeta.Name)
		fmt.Println("-----------------------")
	}
	//c.Status(200).JSON(podsInfo)
	return jsonResponse(c, fiber.StatusOK, "OK", 1, podsInfo)
}

func getDeploymentPods(c *fiber.Ctx) error {
	queryParams := c.Queries()
	// Retrieve specific query parameters
	env := queryParams["env"]
	namespace := queryParams["namespace"]
	deploymentName := queryParams["deploymentName"]

	podsInfo, err := controller.GetDeploymentPods(env, deploymentName, namespace)
	if err != nil {
		fmt.Printf("Error  GetDeploymentPods: %v\n", err)
		// os.Exit(1)
		return jsonResponse(c, fiber.StatusOK, "OK", 1, err)
	}
	for _, pods := range podsInfo {
		fmt.Printf("Name: %s\n", pods.ObjectMeta.Name)
		fmt.Println("-----------------------")
	}
	// return c.Status(200).JSON(podsInfo)
	return jsonResponse(c, fiber.StatusOK, "OK", 1, podsInfo)
}

func deletePods(c *fiber.Ctx) error {
	queryParams := c.Queries()
	// Retrieve specific query parameters
	env := queryParams["env"]
	namespace := queryParams["namespace"]
	podName := queryParams["podName"]
	controller.DeletePod(env, namespace, podName)
	return jsonResponse(c, fiber.StatusOK, "OK", 1, "pod deleted")
}
