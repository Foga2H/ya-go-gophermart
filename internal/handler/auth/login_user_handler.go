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

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserHandler struct {
	service *service.LoginUserService
	token   *auth.Token
	logger  *logger.Logger
}

func NewLoginUserHandler(repository repository.UserRepository, token *auth.Token, logger *logger.Logger) *LoginUserHandler {
	return &LoginUserHandler{
		service: service.NewLoginUserService(repository, logger),
		token:   token,
		logger:  logger,
	}
}

func (h *LoginUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonDecoder := json.NewDecoder(r.Body)
	var request LoginUserRequest
	err := jsonDecoder.Decode(&request)
	if err != nil {
		h.logger.Errorf("json decode error: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), request.Login, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			h.logger.Errorf("user not found: %v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.logger.Errorf("login error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := h.token.GenerateToken(user.Id)
	if err != nil {
		h.logger.Errorf("generate token error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
