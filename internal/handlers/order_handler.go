package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourname/qr-ordering-system/internal/repository"
	"github.com/yourname/qr-ordering-system/internal/services"
	"github.com/yourname/qr-ordering-system/internal/utils"
)

type OrderHandler struct {
	orderService *services.OrderService
	orderRepo    *repository.OrderRepository
}

func NewOrderHandler(orderService *services.OrderService, orderRepo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{orderService: orderService, orderRepo: orderRepo}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var input services.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.orderService.PlaceOrder(input)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")

	order, err := h.orderRepo.GetByID(orderID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "order not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) ConfirmServed(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")

	if err := h.orderService.ConfirmServed(orderID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "served"})
}

func (h *OrderHandler) FlushOrders(w http.ResponseWriter, r *http.Request) {
	if err := h.orderService.FlushAllOrders(); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to flush orders")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "all orders flushed"})
}
