package repository

import (
	"fmt"
	"sync"
	"time"

	"example.com/order-app/api/model"
)

// OrderRepository はインメモリで注文データを保持する。
type OrderRepository struct {
	mu     sync.RWMutex
	orders []model.Order
}

// NewOrderRepository はシードデータを持つ OrderRepository を生成する。
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: seedOrders()}
}

// List は保持している全注文を返す。
func (r *OrderRepository) List() []model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]model.Order, len(r.orders))
	copy(orders, r.orders)
	return orders
}

// FindByID は指定 ID の注文を返す。見つからない場合は ok が false になる。
func (r *OrderRepository) FindByID(id string) (order model.Order, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, o := range r.orders {
		if o.ID == id {
			return o, true
		}
	}
	return model.Order{}, false
}

func seedOrders() []model.Order {
	statuses := make([]string, 0, 24)
	for i := 0; i < 12; i++ {
		statuses = append(statuses, "pending")
	}
	for i := 0; i < 9; i++ {
		statuses = append(statuses, "shipped")
	}
	for i := 0; i < 3; i++ {
		statuses = append(statuses, "canceled")
	}

	base := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	orders := make([]model.Order, len(statuses))
	for i, status := range statuses {
		orders[i] = model.Order{
			ID:        fmt.Sprintf("ord-%03d", i+1),
			Customer:  fmt.Sprintf("customer-%02d", i+1),
			Status:    status,
			Amount:    1000 + i*500,
			OrderedAt: base.Add(time.Duration(i) * 24 * time.Hour).Format(time.RFC3339),
		}
	}
	return orders
}
