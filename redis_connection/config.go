package redis_connection

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sjian_mstr/cluster-management/gaget"
)

var redisClient *redis.Client

// Initialize the Redis client
func initRedisClient(addr string) *redis.Client {
	// Create a new Redis client instance
	client := redis.NewClient(&redis.Options{
		Addr:     addr, // Replace with your Redis server address
		Password: "",   // No password for local Redis instance, update if needed
		DB:       1,    // Default DB
	})

	// Ping the Redis server to check if it's reachable
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}

func init() {
	// get info
	connectionInfo, err := gaget.GetRedisConnectionInfo()
	envTab := "dev"
	addr := connectionInfo[envTab]["HOST"] + ":" + connectionInfo[envTab]["PORT"]
	if err != nil {
		fmt.Println("error values")
		fmt.Println(err)
		// fmt.Errorf("The value must be greater than or equal to 1")
	}
	// Initialize the Redis client
	redisClient = initRedisClient(addr)

	ConsumeMessages("mystream", "0")
}
