package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
)

type QRService struct {
	tableRepo *repository.TableRepository
	baseURL   string
	outputDir string
}

func NewQRService(tableRepo *repository.TableRepository, baseURL string) *QRService {
	return &QRService{
		tableRepo: tableRepo,
		baseURL:   baseURL,
		outputDir: "./web/static/images/qrcodes",
	}
}

// CreateTable creates a table record and generates its QR code image
func (s *QRService) CreateTable(tableNumber int) (*models.Table, string, error) {
	orderURL := fmt.Sprintf("%s/order/%d", s.baseURL, tableNumber)

	table, err := s.tableRepo.Create(tableNumber, orderURL)
	if err != nil {
		return nil, "", fmt.Errorf("creating table record: %w", err)
	}

	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return nil, "", fmt.Errorf("creating output dir: %w", err)
	}

	filename := fmt.Sprintf("table-%d.png", tableNumber)
	fullPath := filepath.Join(s.outputDir, filename)

	if err := qrcode.WriteFile(orderURL, qrcode.Medium, 512, fullPath); err != nil {
		return nil, "", fmt.Errorf("generating QR image: %w", err)
	}

	return table, filename, nil
}

// GenerateBatch creates tables 1..count in one go — useful for initial shop setup
func (s *QRService) GenerateBatch(count int) ([]models.Table, error) {
	var tables []models.Table
	for i := 1; i <= count; i++ {
		table, _, err := s.CreateTable(i)
		if err != nil {
			return nil, fmt.Errorf("table %d: %w", i, err)
		}
		tables = append(tables, *table)
	}
	return tables, nil
}
