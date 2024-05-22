package app

import (
	"fmt"
	"log"

	"github.com/NessibeliY/API/internal/client"
	"github.com/NessibeliY/API/internal/config"
	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/services"
	"github.com/NessibeliY/API/pkg"
	"github.com/gin-gonic/gin"
)

// TODO read cronjobs
// TODO allowed origins корсы? добавить

func Run() {
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

	database := database.NewDatabase(db, rdb)
	services := services.NewServices(database)
	client := client.NewClient(services)

	client.Routes(router)

	// TODO graceful shutdown
	err = router.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
