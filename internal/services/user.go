package services

import (
	"context"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/pkg"
)

type UserStorage interface {
	CreateUser(context.Context, *dto.SignupRequest, []byte) error
	CheckUser(context.Context, *dto.LoginRequest) error
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // TODO transport layer передает аргументос контекст, убрать отсюда
	defer cancel()

	err := us.userStorage.CheckUser(ctx, request) // лучше обработать дто тут
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // TODO почему 3 секунды, почему 20 секунд
	defer cancel()

	// TODO check if user already exists

	err = us.userStorage.CreateUser(ctx, request, hash)
	if err != nil {
		return err
	}

	return nil
}
