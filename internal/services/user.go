package services

import (
	"context"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/pkg"
)

type UserStorage interface {
	CreateUser(context.Context, *dto.SignupRequest, []byte) error
	CheckUserInDB(context.Context, *dto.LoginRequest) error
}

type UserServices struct {
	userStorage UserStorage
}

func NewUserServices(userStorage UserStorage) *UserServices {
	return &UserServices{
		userStorage: userStorage,
	}
}

func (us *UserServices) LoginUser(request *dto.LoginRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := us.userStorage.CheckUserInDB(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServices) SignupUser(request *dto.SignupRequest) error {
	hash, err := pkg.HashPassword(request.Password)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = us.userStorage.CreateUser(ctx, request, hash)
	if err != nil {
		return err
	}

	return nil
}
