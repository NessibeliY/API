package database

import (
	"context"
	"database/sql"
	"time"
)

func Init(db *sql.DB) error {
	query1 := `
	CREATE TABLE IF NOT EXISTS users (
		id uuid PRIMARY KEY,
		username text NOT NULL,
		email text UNIQUE NOT NULL,
		password_hashed bytea NOT NULL 
	);`
	// TODO is it allowed to save with bytea or better string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // TODO move to main
	defer cancel()

	_, err := db.ExecContext(ctx, query1)
	if err != nil {
		return err
	}

	query2 := `
	CREATE TABLE IF NOT EXISTS documents (
		id bigserial PRIMARY KEY,
		title text UNIQUE NOT NULL,
		content text NOT NULL,
		image_path text,
		author_id uuid REFERENCES users(id) ON DELETE CASCADE,
		date_created timestamp NOT NULL,
		date_expired timestamp
	);`

	_, err = db.ExecContext(ctx, query2)
	return err
}
