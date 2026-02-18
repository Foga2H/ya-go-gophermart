package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	repositorymock "github.com/Foga2H/ya-go-gophermart/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUserServiceLoginSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	password := "test-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build test hash: %v", err)
	}

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "test").
		Return(model.User{
			ID:        "u-1",
			Login:     "test",
			Password:  string(hash),
			CreatedAt: time.Now(),
		}, nil)

	svc := NewLoginUserService(repo, logger.NewNopLogger())

	user, err := svc.Login(context.Background(), "test", password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == "" {
		t.Fatal("expected user id")
	}
}

func TestLoginUserServiceLoginUserNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "missing").
		Return(model.User{}, sql.ErrNoRows)

	svc := NewLoginUserService(repo, logger.NewNopLogger())

	_, err := svc.Login(context.Background(), "missing", "password")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestLoginUserServiceLoginGetUserFailed(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "test").
		Return(model.User{}, errors.New("db failure"))

	svc := NewLoginUserService(repo, logger.NewNopLogger())

	_, err := svc.Login(context.Background(), "test", "password")
	if !errors.Is(err, ErrGetUserFailed) {
		t.Fatalf("expected ErrGetUserFailed, got %v", err)
	}
}

func TestLoginUserServiceLoginInvalidPassword(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to build test hash: %v", err)
	}

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "test").
		Return(model.User{
			ID:        "u-2",
			Login:     "test",
			Password:  string(hash),
			CreatedAt: time.Now(),
		}, nil)

	svc := NewLoginUserService(repo, logger.NewNopLogger())

	_, err = svc.Login(context.Background(), "test", "wrong-password")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
