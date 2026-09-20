package repository

import (
	"database/sql"

	"github.com/yourname/qr-ordering-system/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO orders (id, table_id, status, total, payment_method) VALUES ($1, $2, $3, $4, $5)`,
		order.ID, order.TableID, order.Status, order.Total, order.PaymentMethod,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, product_id, quantity, unit_price, notes)
			 VALUES ($1, $2, $3, $4, $5)`,
			order.ID, item.ProductID, item.Quantity, item.UnitPrice, item.Notes,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}


func (r *OrderRepository) GetPending() ([]models.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, table_id, status, total, payment_method, created_at FROM orders
		 WHERE status IN ('pending', 'preparing', 'ready') ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TableID, &o.Status, &o.Total, &o.PaymentMethod, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) GetPendingWithItems() ([]models.Order, error) {
	orders, err := r.GetPending()
	if err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := r.getItemsByOrderID(orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}


func (r *OrderRepository) getItemsByOrderID(orderID string) ([]models.OrderItem, error) {
	rows, err := r.db.Query(
		`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.notes, p.name
		 FROM order_items oi
		 JOIN products p ON p.id = oi.product_id
		 WHERE oi.order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Notes, &item.ProductName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepository) GetByID(orderID string) (*models.Order, error) {
	var o models.Order
	err := r.db.QueryRow(
		`SELECT id, table_id, status, total, payment_method, created_at FROM orders WHERE id = $1`,
		orderID,
	).Scan(&o.ID, &o.TableID, &o.Status, &o.Total, &o.PaymentMethod, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	items, err := r.getItemsByOrderID(o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

func (r *OrderRepository) GetActiveByTableID(tableID int) (*models.Order, error) {
	var o models.Order
	err := r.db.QueryRow(
		`SELECT id, table_id, status, total, payment_method, created_at FROM orders
		 WHERE table_id = $1 AND status IN ('pending', 'preparing', 'ready')
		 ORDER BY created_at DESC LIMIT 1`,
		tableID,
	).Scan(&o.ID, &o.TableID, &o.Status, &o.Total, &o.PaymentMethod, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}


func (r *OrderRepository) UpdateStatus(orderID string, status models.OrderStatus) error {
	result, err := r.db.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, orderID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *OrderRepository) GetHistory(limit int) ([]models.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, table_id, status, total, payment_method, created_at FROM orders
		 WHERE status = 'served' ORDER BY created_at DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TableID, &o.Status, &o.Total, &o.PaymentMethod, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}


type RevenueSummary struct {
	TotalOrders  int     `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
}

func (r *OrderRepository) GetRevenueSummary() (*RevenueSummary, error) {
	var s RevenueSummary
	err := r.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(total), 0) FROM orders WHERE status = 'served'`,
	).Scan(&s.TotalOrders, &s.TotalRevenue)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *OrderRepository) GetRevenueToday() (*RevenueSummary, error) {
	var s RevenueSummary
	err := r.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(total), 0) FROM orders
		 WHERE status = 'served' AND date(created_at) = date('now')`,
	).Scan(&s.TotalOrders, &s.TotalRevenue)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *OrderRepository) FlushAll() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM order_items`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM orders`); err != nil {
		return err
	}

	return tx.Commit()
}
