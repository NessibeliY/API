package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/NessibeliY/API/internal/models"
	"github.com/go-redis/redis/v8"
)

const sessionKeyPrefix = "session:"

func NewRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func SetSessionData(sessionID string, sessionUser models.SessionUserClient, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	err = client.Set(ctx, sessionKeyPrefix+sessionID, jsonData, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetSessionData(client *redis.Client, sessionID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := client.Get(ctx, sessionKeyPrefix+sessionID).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func GetSessionUser(client *redis.Client, sessionID string) (*models.SessionUserClient, error) {
	data, err := GetSessionData(client, sessionID)
	if err != nil {
		return nil, err
	}

	var sessionUser models.SessionUserClient
	err = json.Unmarshal([]byte(data), &sessionUser)
	if err != nil {
		return nil, err
	}

	return &sessionUser, nil
}
