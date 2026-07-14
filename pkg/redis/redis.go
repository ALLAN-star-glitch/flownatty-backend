package redis

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

var client *redis.Client

// Init initializes the Redis client
func Init(addr string) error {
    client = redis.NewClient(&redis.Options{
        Addr: addr,
    })

    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return err
    }

    return nil
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
    return client
}

// Set stores a value with TTL
func Set(key string, value interface{}, ttl time.Duration) error {
    ctx := context.Background()
    return client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value, returns (value, exists, error)
func Get(key string) (string, bool, error) {
    ctx := context.Background()
    val, err := client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return "", false, nil
        }
        return "", false, err
    }
    return val, true, nil
}

// Delete removes a key
func Delete(key string) error {
    ctx := context.Background()
    return client.Del(ctx, key).Err()
}

// HSet stores a hash
func HSet(key string, values map[string]interface{}) error {
    ctx := context.Background()
    return client.HSet(ctx, key, values).Err()
}

// HGetAll retrieves all fields from a hash, returns (values, exists, error)
func HGetAll(key string) (map[string]string, bool, error) {
    ctx := context.Background()
    result, err := client.HGetAll(ctx, key).Result()
    if err != nil {
        return nil, false, err
    }
    if len(result) == 0 {
        return nil, false, nil
    }
    return result, true, nil
}

// Expire sets TTL on a key
func Expire(key string, ttl time.Duration) error {
    ctx := context.Background()
    return client.Expire(ctx, key, ttl).Err()
}

// HGet gets a specific field from a hash, returns (value, exists, error)
func HGet(key, field string) (string, bool, error) {
    ctx := context.Background()
    val, err := client.HGet(ctx, key, field).Result()
    if err != nil {
        if err == redis.Nil {
            return "", false, nil
        }
        return "", false, err
    }
    return val, true, nil
}

// HSetField sets a single field in a hash
func HSetField(key, field string, value interface{}) error {
    ctx := context.Background()
    return client.HSet(ctx, key, field, value).Err()
}