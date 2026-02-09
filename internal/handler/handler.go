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

type TotalCostResponse struct {
	Total int `json:"total"`
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

// Create subscription
// @Summary Create subscription
// @Description Create new subscription for user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Subscription data"
// @Success 201 "Created"
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions [post]

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

// Delete subscription
// @Summary Delete subscription by ID
// @Tags subscriptions
// @Param id query string true "Subscription ID" format(uuid)
// @Success 201 "Deleted"
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions/delete [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := utils.ParseUUIDFromString(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSubscription(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get subscription by ID
// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id query string true "Subscription ID" format(uuid)
// @Success 200 {object} model.Subscription
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions/get [get]
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := utils.ParseUUIDFromString(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetSubscriptionById(ctx, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		http.Error(w, "failed to encode subscriptions", http.StatusInternalServerError)
		return
	}

}

// List subscriptions
// @Summary Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := h.service.List(ctx)
	if err != nil {
		http.Error(w, "invalid list", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, "failed to encode subscriptions", http.StatusInternalServerError)
		return
	}

}

// Calculate total subscriptions cost
// @Summary Calculate total subscriptions cost
// @Description Calculate total cost for subscriptions in period
// @Tags subscriptions
// @Produce json
// @Param userId query string true "User ID" format(uuid)
// @Param serviceName query string false "Service name"
// @Param start query string true "Start month (MM-YYYY)" example(01-2025)
// @Param end query string true "End month (MM-YYYY)" example(12-2025)
// @Success 200 {object} TotalCostResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем query-параметры
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "missing userId", http.StatusBadRequest)
		return
	}
	userID, err := utils.ParseUUIDFromString(userIDStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	serviceName := r.URL.Query().Get("serviceName") 

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	if startStr == "" || endStr == "" {
		http.Error(w, "missing start or end", http.StatusBadRequest)
		return
	}

	startDate, err := utils.ParseMonthYear(startStr)
	if err != nil {
		http.Error(w, "invalid start date", http.StatusBadRequest)
		return
	}

	endDate, err := utils.ParseMonthYear(endStr)
	if err != nil {
		http.Error(w, "invalid end date", http.StatusBadRequest)
		return
	}

	total, err := h.service.CalculateSubscriptionsTotalCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(TotalCostResponse{Total: total}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
