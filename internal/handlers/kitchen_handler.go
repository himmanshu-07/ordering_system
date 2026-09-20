package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourname/qr-ordering-system/internal/kitchen"
	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/utils"
)

type KitchenHandler struct {
	dashboard *kitchen.DashboardService
}

func NewKitchenHandler(dashboard *kitchen.DashboardService) *KitchenHandler {
	return &KitchenHandler{dashboard: dashboard}
}

func (h *KitchenHandler) GetActiveOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.dashboard.GetActiveOrders()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to load orders")
		return
	}
	utils.RespondJSON(w, http.StatusOK, orders)
}

type UpdateStatusInput struct {
	Status models.OrderStatus `json:"status"`
}

func (h *KitchenHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")

	var input UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.dashboard.AdvanceStatus(orderID, input.Status); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
