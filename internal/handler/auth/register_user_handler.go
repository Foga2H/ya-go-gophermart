package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Foga2H/ya-go-gophermart/internal/auth"
	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/repository"
	service "github.com/Foga2H/ya-go-gophermart/internal/service/auth"
)

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type RegisterUserHandler struct {
	service *service.RegisterUserService
	token   *auth.Token
	logger  *logger.Logger
}

func NewRegisterUserHandler(repository repository.UserRepository, authToken *auth.Token, logger *logger.Logger) *RegisterUserHandler {
	return &RegisterUserHandler{
		service: service.NewRegisterUserService(repository, logger),
		token:   authToken,
		logger:  logger,
	}
}

func (h *RegisterUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonDecoder := json.NewDecoder(r.Body)
	var request RegisterUserRequest
	err := jsonDecoder.Decode(&request)
	if err != nil {
		h.logger.Errorf("invalid json body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := h.service.FindByLogin(r.Context(), request.Login)
	if err != nil {
		if errors.Is(err, service.ErrGetUserFailed) || !errors.Is(err, service.ErrUserNotFound) {
			h.logger.Errorf("find user by login error: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	if user.ID != "" {
		h.logger.Errorf("user with id %s already exists", user.ID)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	user, err = h.service.Create(r.Context(), request.Login, request.Password)
	if err != nil {
		h.logger.Errorf("create user error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := h.token.GenerateToken(user.ID)
	if err != nil {
		h.logger.Errorf("generate token error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
