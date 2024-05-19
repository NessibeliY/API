package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/NessibeliY/API/config"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) { // move to pkg or read about infrastructure
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, errors.Wrap(err, "opening sql")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "connection is not established")
	}

	log.Println("Connected to DB")

	return db, nil
}

func OpenRedisDB(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, errors.Wrap(err, "connection to Redis is not established")
	}

	log.Println("Connected to Redis")
	return rdb, nil
}
