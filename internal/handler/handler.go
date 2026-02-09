package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *slog.Logger
}

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartMonth  string  `json:"start_month"` // "07-2025"
	EndMonth    *string `json:"end_month"`   // optional
}

func NewSubscriptionHandler(
	service *service.SubscriptionService,
	logger *slog.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseUUIDFromString(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	startDate, err := utils.ParseMonthYear(req.StartMonth)
	if err != nil {
		http.Error(w, "invalid start_month", http.StatusBadRequest)
		return
	}

	var endDate *time.Time
	if req.EndMonth != nil {
		t, err := utils.ParseMonthYear(*req.EndMonth)
		if err != nil {
			http.Error(w, "invalid end_month", http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.service.CreateSubscription(ctx, sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
