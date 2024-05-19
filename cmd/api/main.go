package main

import (
	"fmt"
	"log"

	"github.com/NessibeliY/API/config"
	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/database/user"
	"github.com/NessibeliY/API/internal/redis"
	"github.com/NessibeliY/API/internal/services"
	"github.com/NessibeliY/API/internal/transport"
	"github.com/NessibeliY/API/pkg"
	"github.com/gin-gonic/gin"
)

// TODO read cronjobs
// TODO allowed origins корсы? добавить

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Println(err, nil)
		return
	}

	// connect to DB
	db, err := pkg.OpenDB(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = database.Init(db)
	if err != nil {
		log.Println(err)
		return
	}

	// Set up Redis DB
	rdb, err := pkg.OpenRedisDB(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	router := gin.Default()

	database := user.NewUserDatabase(db)
	redisDatabase := redis.NewSessionDatabase(rdb)
	services := services.NewUserServices(database, redisDatabase)
	transport := transport.NewTransport(services)

	transport.Routes(router)

	// TODO graceful shutdown
	err = router.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
