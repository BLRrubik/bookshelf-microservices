package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	res := c.client.Get(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return res.Bytes()
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) GenerateKey(prefix, path, query string) string {
	hash := sha256.Sum256([]byte(path + query))

	return fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hash[:])[:16])
}
