package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NessibeliY/API/internal/models"
	"github.com/go-redis/redis/v8"
)

const sessionKeyPrefix = "session:"

type SessionDatabase struct {
	rdb *redis.Client
}

func NewSessionDatabase(rdb *redis.Client) *SessionDatabase {
	return &SessionDatabase{rdb: rdb}
}

func (sdb *SessionDatabase) SetSessionData(ctx context.Context, key string, value models.SessionUserClient, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = sdb.rdb.Set(ctx, key, p, expiration).Err()//key=user, value=userID
	if err != nil {
		return err
	}

	return nil
}

func (sdb *SessionDatabase) GetSessionData(ctx context.Context, key string, dest *models.SessionUserClient) error {
	p, err := sdb.rdb.Get(ctx, key).Result()
	if err != nil {
		// Check if the error is because the key does not exist
		if err == redis.Nil {
			return fmt.Errorf("key does not exist: %v", key)
		}
		return err
	}

	err = json.Unmarshal([]byte(p), dest)
	if err != nil {
		return err
	}

	return nil
}
