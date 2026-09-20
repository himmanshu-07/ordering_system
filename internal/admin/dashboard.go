package admin

import (
	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
)

type DashboardService struct {
	orderRepo *repository.OrderRepository
}

func NewDashboardService(orderRepo *repository.OrderRepository) *DashboardService {
	return &DashboardService{orderRepo: orderRepo}
}

type DashboardData struct {
	Today   *repository.RevenueSummary `json:"today"`
	AllTime *repository.RevenueSummary `json:"all_time"`
	History []models.Order             `json:"history"`
}

func (s *DashboardService) GetDashboard() (*DashboardData, error) {
	today, err := s.orderRepo.GetRevenueToday()
	if err != nil {
		return nil, err
	}
	allTime, err := s.orderRepo.GetRevenueSummary()
	if err != nil {
		return nil, err
	}
	history, err := s.orderRepo.GetHistory(100)
	if err != nil {
		return nil, err
	}
	return &DashboardData{Today: today, AllTime: allTime, History: history}, nil
}
