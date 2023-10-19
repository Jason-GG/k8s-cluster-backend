package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sjian_mstr/cluster-management/database"
	"github.com/sjian_mstr/cluster-management/models"
)

type ClusterInfo struct {
	// This is not the model, more like a serializer
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func CreateResponseUser(clusterInfo models.ClusterInfo) ClusterInfo {
	return ClusterInfo{ID: clusterInfo.ID, Name: clusterInfo.Name}
}

func CreateClusterInfo(c *fiber.Ctx) error {
	var clusterInfo models.ClusterInfo

	if err := c.BodyParser(&clusterInfo); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&clusterInfo)
	responseClusterInfo := CreateResponseUser(clusterInfo)
	return c.Status(200).JSON(responseClusterInfo)
}

func GetClusterInfos(c *fiber.Ctx) error {
	users := []models.ClusterInfo{}
	database.Database.Db.Find(&users)
	responseUsers := []ClusterInfo{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func findClusterInfo(id int, user *models.ClusterInfo) error {
	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func GetClusterInfo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var clusterInfo models.ClusterInfo

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findClusterInfo(id, &clusterInfo); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseClusterInfo := CreateResponseUser(clusterInfo)

	return c.Status(200).JSON(responseClusterInfo)
}
