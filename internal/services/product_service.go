package services

import (
	"errors"

	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

type CreateProductInput struct {
	CategoryID  int     `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

func (s *ProductService) Create(input CreateProductInput) (*models.Product, error) {
	if input.Name == "" {
		return nil, errors.New("product name is required")
	}
	if input.Price <= 0 {
		return nil, errors.New("product price must be greater than zero")
	}
	if input.CategoryID == 0 {
		return nil, errors.New("category_id is required")
	}
	return s.productRepo.Create(input.CategoryID, input.Name, input.Description, input.Price, input.ImageURL)
}

func (s *ProductService) Delete(id int) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) ToggleAvailability(id int) (*models.Product, error) {
	return s.productRepo.ToggleAvailability(id)
}
