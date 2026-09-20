package services

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/yourname/qr-ordering-system/internal/models"
	"github.com/yourname/qr-ordering-system/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo}
}

type OrderItemInput struct {
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Notes     string `json:"notes"`
}

type CreateOrderInput struct {
	TableID       int              `json:"table_id"`
	PaymentMethod string           `json:"payment_method"`
	Items         []OrderItemInput `json:"items"`
}

func (s *OrderService) PlaceOrder(input CreateOrderInput) (*models.Order, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	if input.PaymentMethod != "upi" && input.PaymentMethod != "cash" {
		return nil, errors.New("payment method must be 'upi' or 'cash'")
	}

	existing, err := s.orderRepo.GetActiveByTableID(input.TableID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("failed to check existing orders for table")
	}
	if existing != nil {
		return nil, errors.New("this table already has an active order")
	}

	order := &models.Order{
		ID:            uuid.NewString(),
		TableID:       input.TableID,
		Status:        models.StatusPending,
		PaymentMethod: input.PaymentMethod,
	}

	var total float64
	for _, item := range input.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("invalid product in order")
		}
		if !product.IsAvailable {
			return nil, errors.New(product.Name + " is currently unavailable")
		}

		orderItem := models.OrderItem{
			ProductID: product.ID,
			ProductName: product.Name,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
			Notes:     item.Notes,
		}
		order.Items = append(order.Items, orderItem)
		total += product.Price * float64(item.Quantity)
	}

	order.Total = total

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}


func (s *OrderService) ConfirmServed(orderID string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}
	if order.Status != models.StatusReady {
		return errors.New("order is not ready to be confirmed yet")
	}
	return s.orderRepo.UpdateStatus(orderID, models.StatusServed)
}

func (s *OrderService) FlushAllOrders() error {
	return s.orderRepo.FlushAll()
}
