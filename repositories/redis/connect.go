package redis

import (
	// Go Internal Packages
	"context"
	"log"

	// External Packages
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Connect creates a new RedisHandler instance consisting of master and slave clients
// and returns the instance and error if any and pings the master and slave clients.
func Connect(ctx context.Context, logger *zap.Logger, uri string) (*redis.Client, error) {
	// Configure the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     uri, // Redis server address
		Password: "",  // No password set
		DB:       0,   // Default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return rdb, nil
}
