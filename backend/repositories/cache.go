package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// TTL do Cache
const cacheTTL = 1 * time.Second

// Prefixos das chaves de cache
const (
	TotalCacheKey       = "stats:total:%d"
	ParticipantCacheKey = "stats:participant:%d"
	HourlyCacheKey      = "stats:hourly:%d"
)

// RedisClient é o cliente Redis global
var RedisClient *redis.Client

func InitRedis(host, port string) {
	// Padrão para localhost:6379 se as variáveis de ambiente não estiverem definidas
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Endereço do servidor Redis
		Password: "",        // Sem senha definida
		DB:       0,         // Usa o DB padrão
	})

	// Testa a conexão Redis
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Continuing without cache.", err)
	} else {
		log.Println("Connected to Redis successfully")
	}
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis connection closed")
	}
}

func GetFromCache(ctx context.Context, key string, result interface{}) (bool, error) {
	if RedisClient == nil {
		return false, nil
	}

	data, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Chave não existe no cache
		return false, nil
	} else if err != nil {
		// Erro ao acessar o Redis
		log.Printf("Redis error: %v", err)
		return false, err
	}

	// Desserializa os dados em cache
	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func SetCache(ctx context.Context, key string, data interface{}) error {
	if RedisClient == nil {
		return nil
	}

	// Serializa os dados
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Define no Redis com TTL
	err = RedisClient.Set(ctx, key, jsonData, cacheTTL).Err()
	if err != nil {
		log.Printf("Redis set error: %v", err)
		return err
	}

	return nil
}
