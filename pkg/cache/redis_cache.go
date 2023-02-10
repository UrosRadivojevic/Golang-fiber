package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/urosradivojevic/health/pkg/model"
)

type RedisCacheInterface interface {
	SetMovie(ctx context.Context, value model.Netflix) error
	Get(ctx context.Context, key string) (model.Netflix, error)
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(c *redis.Client) *RedisCache {
	return &RedisCache{
		client: c,
	}
}

func (cache *RedisCache) SetMovie(ctx context.Context, value model.Netflix) error {

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = cache.client.HSet(context.Background(), "movies", value.ID.Hex(), json).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *RedisCache) Get(ctx context.Context, key string) (model.Netflix, error) {

	val, err := cache.client.HGet(ctx, "movies", key).Result()
	if err != nil {
		return model.Netflix{}, err
	}
	movie := model.Netflix{}
	if err := json.Unmarshal([]byte(val), &movie); err != nil {
		return model.Netflix{}, err
	}
	return movie, nil
}

func (cache *RedisCache) Delete(ctx context.Context, key string) error {
	err := cache.client.HDel(ctx, "movies", key).Err()
	if err != nil {
		return err
	}
	return nil
}
