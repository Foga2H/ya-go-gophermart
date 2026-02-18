package handler

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Foga2H/ya-go-gophermart/internal/auth"
	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	repositorymock "github.com/Foga2H/ya-go-gophermart/internal/repository/mock"
	"github.com/golang/mock/gomock"
)

func TestRegisterUserHandlerInvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testLogger := logger.NewNopLogger()
	h := NewRegisterUserHandler(repositorymock.NewMockUserRepository(ctrl), auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestRegisterUserHandlerUserAlreadyExists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "taken").
		Return(model.User{
			ID:        "u-1",
			Login:     "taken",
			Password:  "hash",
			CreatedAt: time.Now(),
		}, nil)

	testLogger := logger.NewNopLogger()
	h := NewRegisterUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(`{"login":"taken","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestRegisterUserHandlerFindByLoginFailure(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "new").
		Return(model.User{}, errors.New("db failure"))

	testLogger := logger.NewNopLogger()
	h := NewRegisterUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(`{"login":"new","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestRegisterUserHandlerCreateFailure(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "new").
		Return(model.User{}, sql.ErrNoRows)
	repo.EXPECT().
		Create(gomock.Any(), "new", gomock.Any()).
		Return(model.User{}, errors.New("insert failure"))

	testLogger := logger.NewNopLogger()
	h := NewRegisterUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(`{"login":"new","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestRegisterUserHandlerCreateSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "new").
		Return(model.User{}, sql.ErrNoRows)
	repo.EXPECT().
		Create(gomock.Any(), "new", gomock.Any()).
		Return(model.User{
			ID:        "u-2",
			Login:     "new",
			Password:  "hashed-password",
			CreatedAt: time.Now(),
		}, nil)

	testLogger := logger.NewNopLogger()
	h := NewRegisterUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(`{"login":"new","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusOK)
	}
}
