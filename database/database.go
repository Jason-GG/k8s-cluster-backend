package database

import (
	"fmt"
	"log"
	"os"

	"github.com/sjian_mstr/cluster-management/gaget"
	"github.com/sjian_mstr/cluster-management/models"

	// "gorm.io/driver/sqlite"
	"encoding/json"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ReadJSONFile(path string) (map[string]map[string]string, error) {
	// Read the JSON file.
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Decode the JSON data into a map of strings to interface{} values.
	var values map[string]map[string]string
	err = json.Unmarshal(jsonData, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func init() {

	sqlConnection, err := gaget.GotMysqlDBConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	// sqlConnection := "admin:milieu0.@tcp(a1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       sqlConnection, // DSN data source name
		DefaultStringSize:         256,           // string
		DisableDatetimePrecision:  true,          // datetime，MySQL 5.6 before, it will not support
		DontSupportRenameIndex:    true,          //
		DontSupportRenameColumn:   true,          //
		SkipInitializeWithVersion: false,         //
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}

	log.Println("Connected Successfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	db.AutoMigrate(&models.ClusterInfo{}, &models.User{})
	Database = DbInstance{
		Db: db,
	}
}
