package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache TTL
const cacheTTL = 60 * time.Second // 1 minute

// Cache key prefixes
const (
	TotalCacheKey       = "stats:total:%d"
	ParticipantCacheKey = "stats:participant:%d"
	HourlyCacheKey      = "stats:hourly:%d"
)

// RedisClient is the global Redis client
var RedisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis(host, port string) {
	// Default to localhost:6379 if environment variables are not set
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Redis server address
		Password: "",        // No password set
		DB:       0,         // Use default DB
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Continuing without cache.", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
}

// CloseRedis closes the Redis client
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis connection closed")
	}
}

// GetFromCache retrieves data from cache
func GetFromCache(ctx context.Context, key string, result interface{}) (bool, error) {
	if RedisClient == nil {
		return false, nil
	}

	data, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist in cache
		return false, nil
	} else if err != nil {
		// Error accessing Redis
		log.Printf("Redis error: %v", err)
		return false, err
	}

	// Unmarshal the cached data
	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return false, err
	}

	return true, nil
}

// SetCache stores data in cache
func SetCache(ctx context.Context, key string, data interface{}) error {
	if RedisClient == nil {
		return nil
	}

	// Marshal the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set in Redis with TTL
	err = RedisClient.Set(ctx, key, jsonData, cacheTTL).Err()
	if err != nil {
		log.Printf("Redis set error: %v", err)
		return err
	}

	return nil
}
