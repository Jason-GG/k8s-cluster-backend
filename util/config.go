package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var SingletonConfig = NewConfig()
var defaultPrefix = "dev"
var ConfigfileName string

type config struct {
	v *viper.Viper
}

type FileInfo struct {
	Name string
	Path string
}

func NewConfig() map[string]FileInfo {
	c, err := defaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func (c *config) setMode(m string) {
	c.v.SetConfigName(m)
	c.v.ReadInConfig()
}

func (c *config) GetRunMode() string {
	return os.Getenv("RUN_ENV")
}

func (c *config) Get(key string) interface{} {
	return c.v.Get(key)
}

func (c *config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *config) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

func (c *config) GetInt32(key string) int32 {
	return c.v.GetInt32(key)
}

func (c *config) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}

func (c *config) GetEnv(env string) string {
	return os.Getenv(env)
}

func (c *config) Unmarshal(raw interface{}) error {
	return c.v.Unmarshal(raw)
}

func (c *config) UnmarshalKey(key string, val interface{}) error {
	return c.v.UnmarshalKey(key, val)
}

func (c *config) watchConfig() error {
	ctx, cancel := context.WithCancel(context.Background())

	c.v.WatchConfig()

	//监听回调函数
	watch := func(e fsnotify.Event) {
		fmt.Printf("Config file is changed: %s \n", e.String())
		cancel()
	}

	c.v.OnConfigChange(watch)
	<-ctx.Done()
	return nil
}

func ListFilesInDirectory(dirPath string) ([]string, error) {
	// Read the directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Create a slice to store file names
	fileNames := make([]string, 0, len(files))

	// Iterate over the files and add their names to the slice
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}

func defaultConfig() (map[string]FileInfo, error) {
	dir, _ := os.Getwd()

	dirPath := filepath.Join(dir, "config")

	// Read the directory
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Create a map to store file information
	fileMap := make(map[string]FileInfo)

	// Iterate over the files and add their information to the map
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			filePath := filepath.Join(dirPath, fileName)
			fileMap[fileName] = FileInfo{
				Name: fileName,
				Path: filePath,
			}
			fmt.Printf("===================>>>>>>>fileMap")
			fmt.Printf(fileName)
		}
	}

	return fileMap, nil
}
