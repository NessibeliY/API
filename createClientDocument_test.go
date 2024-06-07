package document_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NessibeliY/API/internal/client"
	"github.com/NessibeliY/API/internal/config"
	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/services"
	"github.com/NessibeliY/API/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupServer() *gin.Engine {
	cfg, err := config.Load()
	if err != nil {
		log.Println(err, nil)
		return nil
	}

	// connect to DB
	db, err := pkg.OpenDB(cfg)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	err = database.Init(db)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Set up Redis DB
	rdb, err := pkg.OpenRedisDB(cfg)
	if err != nil {
		log.Println(err)
		return nil
	}

	router := gin.Default()

	database := database.NewDatabase(db, rdb)
	services := services.NewServices(database)
	client := client.NewClient(services)

	client.Routes(router)

	err = router.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}

	return router
}

func runTestServer() *httptest.Server {
	return httptest.NewServer(setupServer())
}

func TestCreateClientDocument(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

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
