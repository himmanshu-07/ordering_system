package handlers

import (
	"net/http"

	"github.com/yourname/qr-ordering-system/internal/utils"
)

type ConfigHandler struct {
	upiID          string
	whatsappNumber string
}

func NewConfigHandler(upiID, whatsappNumber string) *ConfigHandler {
	return &ConfigHandler{upiID: upiID, whatsappNumber: whatsappNumber}
}

func (h *ConfigHandler) GetPublicConfig(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"upi_id":          h.upiID,
		"whatsapp_number": h.whatsappNumber,
	})
}
