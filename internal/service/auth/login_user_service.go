package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	"github.com/Foga2H/ya-go-gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserService struct {
	repository repository.UserRepository
	logger     *logger.Logger
}

func NewLoginUserService(userRepository repository.UserRepository, logger *logger.Logger) *LoginUserService {
	return &LoginUserService{
		repository: userRepository,
		logger:     logger,
	}
}

func (s *LoginUserService) Login(ctx context.Context, login string, password string) (model.User, error) {
	user, err := s.repository.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, ErrGetUserFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, ErrUserNotFound
	}

	return user, nil
}
