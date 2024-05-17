package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserDatabase struct {
	db *sql.DB
}

func NewUserDatabase(db *sql.DB) *UserDatabase {
	return &UserDatabase{db: db}
}

func (udb *UserDatabase) CheckUserInDB(ctx context.Context, request *dto.LoginRequest) error {
	query := `SELECT password_hashed FROM users
	WHERE email=$1`

	var storedPassword []byte
	err := udb.db.QueryRowContext(ctx, query, request.Email).Scan(&storedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(request.Password))
	if err != nil {
		return models.ErrPasswordNotCorrect
	}

	return nil
}

func (udb *UserDatabase) CreateUser(ctx context.Context, request *dto.SignupRequest, hashedPassword []byte) error {
	query := `INSERT INTO users (id, username, email, password_hashed)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	user := models.User{
		ID:       uuid.New(),
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	err := udb.db.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, user.Password).Scan(&user.ID)
	switch {
	case errors.Is(err, errors.New(`pq: duplicate key value violates unique constraint "users_email_key"`)):
		return models.ErrDuplicateEmail
	case err != nil:
		return err
	}

	return nil
}
