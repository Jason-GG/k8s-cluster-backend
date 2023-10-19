package redis_connection

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Set a key-value pair in Redis
func SetKey(key, value string) error {
	return redisClient.Set(context.Background(), key, value, 0).Err()
}

// Get the value of a key from Redis
func GetKey(key string) (string, error) {
	return redisClient.Get(context.Background(), key).Result()
}

func ProduceMessages(streamKey string, logs []string) {
	// Publish messages to a Redis stream
	for _, log := range logs {
		_, err := redisClient.XAdd(context.Background(), &redis.XAddArgs{
			Stream: streamKey,
			Values: map[string]interface{}{"message": log},
		}).Result()
		if err != nil {
			fmt.Println("Error publishing message:", err)
		}
	}
}

func ConsumeLatestMessage(streamName string) {
	for {
		result, err := redisClient.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{streamName, "2"},
			Count:   1, // Number of messages to read per call
		}).Result()
		if err != nil {
			fmt.Println("Error reading from stream:", err)
			break
		}

		for _, message := range result {
			for _, xMessage := range message.Messages {
				messageValue, ok := xMessage.Values["message"].(string)
				if ok {
					fmt.Printf("Received latest message: %s\n", messageValue)
				}
			}
		}
	}
}

func ConsumeMessages(streamName, sId string) []map[string]string {
	// Minimum and maximum message IDs (inclusive)
	var streamKeyValue []map[string]string
	minID := "-"
	maxID := "+"
	result, err := redisClient.XRange(context.Background(), streamName, minID, maxID).Result()
	if err != nil {
		fmt.Println("Error retrieving stream keys:", err)
		return nil
	}
	count := 0
	for _, message := range result {
		key := message.ID
		value, ok := message.Values["message"].(string)
		if ok {
			count += 1
			fmt.Printf("Key: %s, Value: %s\n, Count: %d\n", key, value, count)
			source := make(map[string]string)
			source[key] = value
			streamKeyValue = append(streamKeyValue, source)
		}
	}

	return streamKeyValue

}
