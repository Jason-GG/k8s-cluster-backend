package gaget

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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

func GotMysqlDBConnection() (string, error) {
	// db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})
	dir, _ := os.Getwd()
	var filePath = filepath.Join(dir, "gaget/sql_info.json") //"./gaget/sql_info.json"
	// var envValue = os.Getenv("ENV_EXE")
	// envValue := "dev"
	envValue := "prod"
	data, err := ReadJSONFile(os.ExpandEnv(filePath))
	if err != nil {
		fmt.Println("error values")
		fmt.Println(err)
		return "", fmt.Errorf("The value must be greater than or equal to 1")
	}

	sqlConnection := data[envValue]["MYSQL_USER"] + ":" + data[envValue]["MYSQL_PASSWORD"] + "@tcp(" +
		data[envValue]["MYSQL_HOST"] + ":3306)/" + data[envValue]["MYSQL_DB"] + "?charset=utf8&parseTime=True&loc=Local"
	return sqlConnection, nil

}

func GetRedisConnectionInfo() (map[string]map[string]string, error) {
	// db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})
	dir, _ := os.Getwd()
	var filePath = filepath.Join(dir, "gaget/redis_info.json")
	fmt.Printf("=====================++>>>>>>filePath", filePath)
	data, err := ReadJSONFile(os.ExpandEnv(filePath))
	if err != nil {
		fmt.Println("error values")
		fmt.Println(err)
		return nil, fmt.Errorf("The value must be greater than or equal to 1")
	}

	return data, nil
}
