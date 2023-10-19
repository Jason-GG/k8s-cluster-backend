package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/controller"
)

func getNamespaceListInfo(c *fiber.Ctx) error {
	cluster := c.Query("cluster")
	namesapceInfo, err := controller.ListNamespaces(cluster)
	if err != nil {
		fmt.Printf("Error lgetNamespaceListInfo: %v\n", err)
		// os.Exit(1)
		return c.Status(200).JSON("Error getNamespaceListInfo")
	}
	var newArray []string
	for _, ns := range namesapceInfo {
		newArray = append(newArray, ns.ObjectMeta.Name)
	}

	response := struct {
		Data []string `json:"data"`
		Msg  string   `json:"msg"`
		Code int      `json:"code"`
	}{
		Data: newArray,
		Msg:  "ok",
		Code: 1,
	}
	return c.Status(200).JSON(response)
}

func getNamespaceListDetailInfo(c *fiber.Ctx) error {
	cluster := c.Query("cluster")
	namesapceInfo, err := controller.ListNamespaces(cluster)
	if err != nil {
		fmt.Printf("Error lgetNamespaceListInfo: %v\n", err)
		// os.Exit(1)
		return c.Status(200).JSON("Error getNamespaceListInfo")
	}

	return jsonResponse(c, fiber.StatusOK, "OK", 1, namesapceInfo)
}

func verifyNamespaceInfo(c *fiber.Ctx) error {
	env := c.Query("env")
	cluster := "microstrategy"
	namesapceInfo, err := controller.ListNamespaces(cluster)
	if err != nil {
		fmt.Printf("Error lgetNamespaceListInfo: %v\n", err)
		// os.Exit(1)
		return c.Status(200).JSON("Error getNamespaceListInfo")
	}
	for _, ns := range namesapceInfo {
		fmt.Println("judge string :", env, ns.ObjectMeta.Name)
		if strings.Contains(ns.ObjectMeta.Name, env) {
			return c.Status(200).JSON("true")
		}
	}

	return c.Status(200).JSON("false")
}
