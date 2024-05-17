package redis

import (
	"context"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/go-redis/redis/v8"
)

func SetCacheInRedis() (*redis.Client, error) {
	store := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := store.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		MaxAge: int(1 * time.Minute / time.Second),
	})

	return store, nil
}
