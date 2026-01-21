package redisfactory

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	initialDelayEnvKey = "REDIS_CONNECTION_INITIAL_DELAY"
	maxDelayEnvKey     = "REDIS_CONNECTION_MAX_DELAY"
)

func FromURL(ctx context.Context, maxTries int, url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return FromOptions(ctx, maxTries, opts)
}

func FromAddress(ctx context.Context, maxTries int, addr string, password string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
	}
	return FromOptions(ctx, maxTries, opts)
}

func FromOptions(ctx context.Context, maxTries int, opts *redis.Options) (*redis.Client, error) {
	redisClient := redis.NewClient(opts)
	tries := 0

	// how many ms to start out the connection delay
	delay := initialDelay()
	for tries < maxTries {
		tries++
		_, err := redisClient.Ping(ctx).Result() // ignoring the "PONG" reply
		if err != nil {
			log.Printf("Unable to connect to Redis, attempt %d: %v", tries, err)
			delay = getTimeout(delay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		} else {
			break
		}
	}
	if tries == maxTries {
		return nil, fmt.Errorf("could not connect to Redis after %d tries", tries)
	}
	return redisClient, nil
}

// gets a delay that has some randomness to help prevent "thundering herd" connections
func getTimeout(prior int) int {
	delay := prior + rand.Intn(prior/2)
	if delay >= maxDelay() {
		return maxDelay() + rand.Intn(100) - 50
	}
	return delay
}

func initialDelay() int {
	return envIntOrDefault(initialDelayEnvKey, 1000)
}

func maxDelay() int {
	return envIntOrDefault(maxDelayEnvKey, 15000)
}

func envIntOrDefault(key string, defaultValue int) int {
	envVal := os.Getenv(key)
	if envVal == "" {
		return defaultValue
	}
	r, err := strconv.Atoi(envVal)
	if err != nil {
		return defaultValue
	}
	return r
}
