package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourname/qr-ordering-system/internal/services"
	"github.com/yourname/qr-ordering-system/internal/utils"
)

type TableHandler struct {
	tableService *services.TableService
}

func NewTableHandler(tableService *services.TableService) *TableHandler {
	return &TableHandler{tableService: tableService}
}

func (h *TableHandler) GetNextAvailable(w http.ResponseWriter, r *http.Request) {
	tableNum, err := h.tableService.GetNextAvailableTable()
	if err != nil {
		if err == services.ErrNoTablesAvailable {
			utils.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]int{"table_number": tableNum})
}

func (h *TableHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	total, err := h.tableService.GetTotalTables()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to load settings")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]int{"total_tables": total})
}

type UpdateTableCountInput struct {
	TotalTables int `json:"total_tables"`
}

func (h *TableHandler) UpdateTableCount(w http.ResponseWriter, r *http.Request) {
	var input UpdateTableCountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.tableService.SetTotalTables(input.TotalTables); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
