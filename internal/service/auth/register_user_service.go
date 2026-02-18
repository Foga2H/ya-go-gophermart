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

var ErrUserCreationFailed = errors.New("user creation failed")
var ErrGetUserFailed = errors.New("get user failed")
var ErrUserNotFound = errors.New("user not found")
var ErrGeneratePasswordHashFailed = errors.New("generate password hash failed")

type RegisterUserService struct {
	repository repository.UserRepository
	logger     *logger.Logger
}

func NewRegisterUserService(userRepository repository.UserRepository, logger *logger.Logger) *RegisterUserService {
	return &RegisterUserService{
		repository: userRepository,
		logger:     logger,
	}
}

func (s *RegisterUserService) Create(ctx context.Context, login string, password string) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, ErrGeneratePasswordHashFailed
	}

	user, err := s.repository.Create(ctx, login, string(hashedPassword))
	if err != nil {
		s.logger.Errorf("Error creating user: %s", err)
		return model.User{}, ErrUserCreationFailed
	}

	return user, nil
}

func (s *RegisterUserService) FindByLogin(ctx context.Context, login string) (model.User, error) {
	user, err := s.repository.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, ErrGetUserFailed
	}

	return user, nil
}
