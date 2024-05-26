package user

import (
	"context"
	"database/sql"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserDatabase struct {
	db *sql.DB
}

func NewUserDatabase(db *sql.DB) *UserDatabase {
	return &UserDatabase{db: db}
}

func (udb *UserDatabase) CreateUser(ctx context.Context, request *dto.SignupRequest, hashedPassword []byte) error { // TODO input only *models.User, without DOT and hashedpassword
	query := `
	INSERT INTO users (id, username, email, password_hashed)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	user := models.User{
		ID:       uuid.New(), // как избежать коллизии uuid, коллизии в map
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	}

	err := udb.db.QueryRowContext(ctx, query, user.ID, user.Username, user.Email, user.Password).Scan(&user.ID)

	// Use pq error code to handle specific PostgreSQL errors
	if pqErr, ok := err.(*pq.Error); ok { // TODO должен быть в сервис, вынести switch
		switch pqErr.Code {
		case "23505": // unique_violation
			if pqErr.Constraint == "users_email_key" {
				return models.ErrDuplicateEmail
			}
		default:
			return err
		}
	}

	return nil
}

func (udb *UserDatabase) CheckUser(ctx context.Context, request *dto.LoginRequest) error {
	query := `SELECT password_hashed FROM users WHERE email=$1;`
	// TODO sql injection

	var storedPassword []byte
	err := udb.db.QueryRowContext(ctx, query, request.Email).Scan(&storedPassword) // TODO QueryRowContext and QueryRow
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(request.Password)) // TODO проверка должна быть на уровень выше
	if err != nil {
		return models.ErrPasswordNotCorrect // TODO errors store in value/errors client or non-client errors
	}

	return nil
}
