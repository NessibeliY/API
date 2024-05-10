package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	// router.SetTrustedProxies([]string{"192.168.1.2"})
	router.SetTrustedProxies(nil)

	err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	// router.POST("/signup", signup)
	router.POST("/login", login)

	// Run the server
	if err = router.Run(":8081"); err != nil {
		log.Fatal("server failed to start:", err)
	}
}

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func login(c *gin.Context) {
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists in DB and if password is correct
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", json.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if json.Password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func ConnectDB() error {
	var err error
	// Connect to DB
	db, err = sql.Open("postgres", "postgres://postgres:12345678@localhost/mydatabase?sslmode=disable")
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	// defer db.Close()

	// Create users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	// Insert "new_user" user to users table
	// username := "new_user"
	// password := "hashed_password"
	// insertStmt := `
	// INSERT INTO users (username, password)
	// VALUES ($1, $2)
	// `
	// _, err = db.Exec(insertStmt, username, password)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("User inserted successfully.")

	// show all users
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	var id int
	var userUsername string
	for rows.Next() {
		err := rows.Scan(&id, &userUsername)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d, Username: %s\n", id, userUsername)
	}

	fmt.Println("All users retrieved.")
	return nil
}
