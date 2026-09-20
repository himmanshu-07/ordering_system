package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yourname/qr-ordering-system/internal/repository"
	"github.com/yourname/qr-ordering-system/internal/services"
	"github.com/yourname/qr-ordering-system/internal/utils"
)

type ProductHandler struct {
	productRepo    *repository.ProductRepository
	categoryRepo   *repository.CategoryRepository
	productService *services.ProductService
}

func NewProductHandler(pr *repository.ProductRepository, cr *repository.CategoryRepository, ps *services.ProductService) *ProductHandler {
	return &ProductHandler{productRepo: pr, categoryRepo: cr, productService: ps}
}

// GetMenu returns categories with their products - the full menu for the QR page
func (h *ProductHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryRepo.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to load menu")
		return
	}

	type CategoryWithProducts struct {
		Category interface{} `json:"category"`
		Products interface{} `json:"products"`
	}

	var menu []CategoryWithProducts
	for _, cat := range categories {
		products, err := h.productRepo.GetAvailableByCategory(cat.ID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "failed to load products")
			return
		}
		menu = append(menu, CategoryWithProducts{Category: cat, Products: products})
	}

	utils.RespondJSON(w, http.StatusOK, menu)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.productRepo.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "product not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input services.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.productService.Create(input)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.productService.Delete(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, "product not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to load products")
		return
	}
	utils.RespondJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) ToggleAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.productService.ToggleAvailability(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "product not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, product)
}
