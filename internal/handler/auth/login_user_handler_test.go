package handler

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Foga2H/ya-go-gophermart/internal/auth"
	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	repositorymock "github.com/Foga2H/ya-go-gophermart/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUserHandlerInvalidJSON(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testLogger := logger.NewNopLogger()
	h := NewLoginUserHandler(repositorymock.NewMockUserRepository(ctrl), auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLoginUserHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "missing").
		Return(model.User{}, sql.ErrNoRows)

	testLogger := logger.NewNopLogger()
	h := NewLoginUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(`{"login":"missing","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestLoginUserHandlerInternalServerError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositorymock.NewMockUserRepository(ctrl)
	repo.EXPECT().
		FindByLogin(gomock.Any(), "test").
		Return(model.User{}, errors.New("db failure"))

	testLogger := logger.NewNopLogger()
	h := NewLoginUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(`{"login":"test","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestLoginUserHandlerSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	password := "pass"
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

	testLogger := logger.NewNopLogger()
	h := NewLoginUserHandler(repo, auth.NewToken("secret", testLogger), testLogger)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(`{"login":"test","password":"pass"}`))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status mismatch: got %d, want %d", rec.Code, http.StatusOK)
	}

	authHeader := rec.Header().Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") == "" {
		t.Fatalf("expected Authorization Bearer token, got %q", authHeader)
	}
}
