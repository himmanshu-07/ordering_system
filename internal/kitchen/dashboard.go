package kitchen

import (
	"errors"

	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
)

type DashboardService struct {
	orderRepo *repository.OrderRepository
}

func NewDashboardService(orderRepo *repository.OrderRepository) *DashboardService {
	return &DashboardService{orderRepo: orderRepo}
}

func (s *DashboardService) GetActiveOrders() ([]models.Order, error) {
	return s.orderRepo.GetPendingWithItems()
}

var validTransitions = map[models.OrderStatus][]models.OrderStatus{
	models.StatusPending:   {models.StatusPreparing},
	models.StatusPreparing: {models.StatusReady},
	models.StatusReady:     {models.StatusServed},
}

func (s *DashboardService) AdvanceStatus(orderID string, newStatus models.OrderStatus) error {
	// In a real system you'd fetch current status first to validate the transition.
	// Keeping it simple: just validate newStatus is one of the known states.
	switch newStatus {
	case models.StatusPreparing, models.StatusReady, models.StatusServed:
		return s.orderRepo.UpdateStatus(orderID, newStatus)
	default:
		return errors.New("invalid status transition")
	}
}
