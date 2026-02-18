package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/middleware"
	"github.com/Foga2H/ya-go-gophermart/internal/repository"
	service "github.com/Foga2H/ya-go-gophermart/internal/service/order"
)

type CreateOrderHandler struct {
	service *service.CreateOrderService
	logger  *logger.Logger
}

func NewCreateOrderHandler(repository repository.OrderRepository, logger *logger.Logger) *CreateOrderHandler {
	return &CreateOrderHandler{
		service: service.NewCreateOrderService(repository, logger),
		logger:  logger,
	}
}

func (h *CreateOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.logger.Error("UserID from context is missing")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("read body err: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(bodyBytes))
	if !h.isValidLuhn(orderNumber) {
		h.logger.Errorf("invalid order number: %q", orderNumber)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	orderFound, err := h.service.FindByOrderID(context.Background(), orderNumber)
	if err != nil {
		if !errors.Is(err, service.ErrOrderNotFound) {
			h.logger.Errorf("find order err: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	if orderFound.Id != "" && orderFound.UserId != userID {
		h.logger.Errorf("order was created by different user: %q", orderNumber)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if orderFound.UserId == userID {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, err = h.service.Create(context.Background(), userID, orderNumber)
	if err != nil {
		h.logger.Errorf("create order err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *CreateOrderHandler) isValidLuhn(number string) bool {
	if number == "" {
		return false
	}

	sum := 0
	shouldDouble := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := number[i]
		if digit < '0' || digit > '9' {
			return false
		}

		value := int(digit - '0')
		if shouldDouble {
			value *= 2
			if value > 9 {
				value -= 9
			}
		}

		sum += value
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}
