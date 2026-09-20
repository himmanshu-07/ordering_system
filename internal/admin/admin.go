package admin

import (
	"encoding/json"
	"net/http"

	"github.com/yourname/qr-ordering-system/internal/services"
	"github.com/yourname/qr-ordering-system/internal/utils"
)

type AdminHandler struct {
	qrService        *services.QRService
	dashboardService *DashboardService
}

func NewAdminHandler(qrService *services.QRService, dashboardService *DashboardService) *AdminHandler {
	return &AdminHandler{qrService: qrService, dashboardService: dashboardService}
}

func (h *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.dashboardService.GetDashboard()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to load dashboard")
		return
	}
	utils.RespondJSON(w, http.StatusOK, data)
}

type CreateTableInput struct {
	TableNumber int `json:"table_number"`
}

func (h *AdminHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	var input CreateTableInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	table, filename, err := h.qrService.CreateTable(input.TableNumber)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"table":    table,
		"qr_image": "/static/images/qrcodes/" + filename,
	})
}

type GenerateBatchInput struct {
	Count int `json:"count"`
}

func (h *AdminHandler) GenerateBatch(w http.ResponseWriter, r *http.Request) {
	var input GenerateBatchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tables, err := h.qrService.GenerateBatch(input.Count)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, tables)
}

