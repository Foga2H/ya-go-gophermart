package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/middleware"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	"github.com/Foga2H/ya-go-gophermart/internal/repository"
	service "github.com/Foga2H/ya-go-gophermart/internal/service/order"
)

type GetUserOrdersHandler struct {
	service *service.GetUserOrdersService
	logger  *logger.Logger
}

func NewGetUserOrdersHandler(repository repository.OrderRepository, logger *logger.Logger) *GetUserOrdersHandler {
	return &GetUserOrdersHandler{
		service: service.NewGetUserOrdersService(repository, logger),
		logger:  logger,
	}
}

type GetUserOrdersResponse struct {
	Number     string            `json:"number"`
	Status     model.OrderStatus `json:"status"`
	Accrual    float32           `json:"accrual,omitempty"`
	UploadedAt string            `json:"uploaded_at"`
}

func (h *GetUserOrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.logger.Error("UserID from context is missing")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	orders, err := h.service.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("get orders err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response []GetUserOrdersResponse
	for _, order := range orders {
		item := GetUserOrdersResponse{
			Number:     order.Number,
			Status:     order.Status,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		}

		if item.Status == model.OrderStatusProcessed {
			item.Accrual = 0
		}

		response = append(response, item)
	}

	resp, err := json.Marshal(response)
	if err != nil {
		h.logger.Errorf("json marshal err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
