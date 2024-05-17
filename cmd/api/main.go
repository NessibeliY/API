package main

import (
	"fmt"
	"log"

	"github.com/NessibeliY/API/config"
	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/services"
	"github.com/NessibeliY/API/internal/transport"

	"github.com/NessibeliY/API/internal/database/user"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Println(err, nil) // return should be after this line?
	}

	// connect to DB
	db, err := openDB(*cfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	router := gin.Default()

	err = database.Init(db)
	if err != nil {
		log.Println(err)
		return
	}

	database := user.NewUserDatabase(db)
	services := services.NewUserServices(database)
	transport := transport.NewTransport(services)

	transport.Routes(router, cfg)

	err = router.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
