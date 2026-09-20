package services

import (
	"database/sql"
	"errors"

	"github.com/yourname/qr-ordering-system/internal/repository"
)

type TableService struct {
	orderRepo    *repository.OrderRepository
	settingsRepo *repository.SettingsRepository
}

func NewTableService(orderRepo *repository.OrderRepository, settingsRepo *repository.SettingsRepository) *TableService {
	return &TableService{orderRepo: orderRepo, settingsRepo: settingsRepo}
}

var ErrNoTablesAvailable = errors.New("we're not accepting new orders right now — all tables are currently occupied")

// GetNextAvailableTable scans table numbers 1..totalTables and returns the first
// one with no active (pending/preparing/ready) order.
func (s *TableService) GetNextAvailableTable() (int, error) {
	total, err := s.settingsRepo.GetTotalTables()
	if err != nil {
		return 0, errors.New("failed to load table configuration")
	}

	for tableNum := 1; tableNum <= total; tableNum++ {
		_, err := s.orderRepo.GetActiveByTableID(tableNum)
		if errors.Is(err, sql.ErrNoRows) {
			return tableNum, nil // no active order on this table — it's free
		}
		if err != nil {
			return 0, errors.New("failed to check table availability")
		}
		// err == nil means an active order exists on this table — keep looking
	}

	return 0, ErrNoTablesAvailable
}

func (s *TableService) GetTotalTables() (int, error) {
	return s.settingsRepo.GetTotalTables()
}

func (s *TableService) SetTotalTables(count int) error {
	if count < 5 || count > 20 {
		return errors.New("total tables must be between 5 and 20")
	}
	return s.settingsRepo.SetTotalTables(count)
}
