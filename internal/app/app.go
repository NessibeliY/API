package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NessibeliY/API/internal/client"
	"github.com/NessibeliY/API/internal/config"
	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/services"
	"github.com/NessibeliY/API/pkg"
	"github.com/gin-gonic/gin"
)

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

	// Graceful shutdown

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server is running on port %v", cfg.Port)

	// Wait for a termination signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful server shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err) // TODO move graceful shutdown to pkg and input only ctx from main.go
	}

	log.Println("Server exiting")
}
