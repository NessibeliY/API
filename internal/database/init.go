package database

import (
	"context"
	"database/sql"
	"time"
)

func Init(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id uuid PRIMARY KEY,
		username text NOT NULL,
		email text UNIQUE NOT NULL,
		password_hashed bytea NOT NULL 
	);`
	// TODO is it allowed to save with bytea or better string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // TODO move to main
	defer cancel()

	_, err := db.ExecContext(ctx, query)
	return err
}
