package document_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NessibeliY/API/internal/client"
	"github.com/NessibeliY/API/internal/config"
	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/services"
	"github.com/NessibeliY/API/pkg"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func createDB(cfg *config.Config) (*sql.DB, error) {
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, "test", "test", "test")

	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, errors.Wrap(err, "opening test sql")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "connection is not established")
	}

	log.Println("Connected to DB")

	return db, nil
}

func dropDB(db *sql.DB, cfg *config.Config) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName)

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop the database: %v", err)
	}

	return nil
}

func setupServer() (*gin.Engine, *sql.DB, *config.Config) {
	cfg, err := config.Load()
	if err != nil {
		log.Println("error loading configs", err)
		return nil, nil, nil
	}

	// create and connect to test DB
	db, err := createDB(cfg)
	if err != nil {
		log.Println(err)
		return nil, nil, nil
	}

	err = database.Init(db)
	if err != nil {
		log.Println(err)
		return nil, nil, nil
	}

	// Set up Redis DB
	rdb, err := pkg.OpenRedisDB(cfg)
	if err != nil {
		log.Println(err)
		return nil, nil, nil
	}

	router := gin.Default()

	database := database.NewDatabase(db, rdb)
	services := services.NewServices(database)
	client := client.NewClient(services)

	client.Routes(router)

	return router, db, cfg
}

func runTestServer() (*httptest.Server, *sql.DB, *config.Config) {
	router, db, cfg := setupServer()
	if router == nil {
		fmt.Println("error is here***************")
		return nil, nil, nil
	}
	return httptest.NewServer(router), db, cfg
}

func TestCreateClientDocument(t *testing.T) {
	ts, db, cfg := runTestServer()
	// if ts == nil || db == nil || cfg == nil {
	// 	t.Fatalf("Failed to set up test server")
	// }
	if ts == nil {
		t.Fatalf("Failed to set up test server")
	}
	if db == nil {
		t.Fatalf("Failed to set up db")
	}
	if cfg == nil {
		t.Fatalf("Failed to set up cfg")
	}
	defer ts.Close()
	defer func() {
		err := dropDB(db, cfg)
		if err != nil {
			log.Fatalf("Failed to drop test database: %v", err)
		}
	}()

	t.Run("it should return 200 when health is ok", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/health", ts.URL))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return validation error when request misses required parameters", func(t *testing.T) {
		// Prepare request body
		requestBody := []byte(`{"content": "Document content"}`)

		// Send POST request to create document endpoint
		resp, err := http.Post(fmt.Sprintf("%s/create-document", ts.URL), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Assert response status code
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("it should return ok when insert new document successfully", func(t *testing.T) {
		// Prepare request body
		requestBody := []byte(`{"title": "Test Document", "content": "Document content"}`)

		// Send POST request to create document endpoint
		resp, err := http.Post(fmt.Sprintf("%s/create-document", ts.URL), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Assert response status code
		assert.Equal(t, 200, resp.StatusCode)
	})
}
