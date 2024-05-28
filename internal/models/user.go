package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrPasswordNotCorrect = errors.New("password not correct")
)

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password []byte
}

type SessionUserClient struct {
	Email         string
	Authenticated bool
}
