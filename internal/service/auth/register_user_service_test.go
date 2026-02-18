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

func TestRegisterUserServiceCreateSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const login = "test-login"
	const plainPassword = "test-password"

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		Create(gomock.Any(), login, gomock.Any()).
		DoAndReturn(func(_ context.Context, gotLogin, gotHashedPassword string) (model.User, error) {
			if gotLogin != login {
				t.Fatalf("login mismatch: got %q, want %q", gotLogin, login)
			}

			if gotHashedPassword == plainPassword {
				t.Fatal("password was not hashed")
			}

			if err := bcrypt.CompareHashAndPassword([]byte(gotHashedPassword), []byte(plainPassword)); err != nil {
				t.Fatalf("hash does not match source password: %v", err)
			}

			return model.User{
				Id:        "u-1",
				Login:     gotLogin,
				Password:  gotHashedPassword,
				CreatedAt: time.Now(),
			}, nil
		})

	svc := NewRegisterUserService(repo, logger.NewNopLogger())

	user, err := svc.Create(context.Background(), login, plainPassword)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Id == "" {
		t.Fatal("expected created user id")
	}
}

func TestRegisterUserServiceCreateRepositoryError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		Create(gomock.Any(), "test-login", gomock.Any()).
		Return(model.User{}, errors.New("insert failed"))

	svc := NewRegisterUserService(repo, logger.NewNopLogger())

	_, err := svc.Create(context.Background(), "test-login", "test-password")
	if !errors.Is(err, ErrUserCreationFailed) {
		t.Fatalf("expected ErrUserCreationFailed, got %v", err)
	}
}

func TestRegisterUserServiceFindByLoginNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "missing-login").
		Return(model.User{}, sql.ErrNoRows)

	svc := NewRegisterUserService(repo, logger.NewNopLogger())

	_, err := svc.FindByLogin(context.Background(), "missing-login")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestRegisterUserServiceFindByLoginRepositoryError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "test-login").
		Return(model.User{}, errors.New("db unavailable"))

	svc := NewRegisterUserService(repo, logger.NewNopLogger())

	_, err := svc.FindByLogin(context.Background(), "test-login")
	if !errors.Is(err, ErrGetUserFailed) {
		t.Fatalf("expected ErrGetUserFailed, got %v", err)
	}
}

func TestRegisterUserServiceFindByLoginSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := model.User{
		Id:        "u-42",
		Login:     "login",
		Password:  "hash",
		CreatedAt: time.Now(),
	}

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "login").
		Return(want, nil)

	svc := NewRegisterUserService(repo, logger.NewNopLogger())

	got, err := svc.FindByLogin(context.Background(), "login")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Id != want.Id || got.Login != want.Login || got.Password != want.Password {
		t.Fatalf("user mismatch: got %+v, want %+v", got, want)
	}
}
