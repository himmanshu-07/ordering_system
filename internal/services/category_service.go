package services

import (
	"errors"

	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

type CreateCategoryInput struct {
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
}

func (s *CategoryService) Create(input CreateCategoryInput) (*models.Category, error) {
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}
	return s.categoryRepo.Create(input.Name, input.DisplayOrder)
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}
