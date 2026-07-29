package service

import (
	"example.com/order-app/api/model"
	"example.com/order-app/api/repository"
)

// OrderService は注文に関するビジネスロジックを担う。
type OrderService struct {
	repo *repository.OrderRepository
}

// NewOrderService は OrderService を生成する。
func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// ListOrders は status/limit/offset に基づき注文一覧と該当件数を返す。
// status が空文字の場合は絞り込みを行わない。
func (s *OrderService) ListOrders(status string, limit, offset int) ([]model.Order, int) {
	orders := s.repo.List()

	if status != "" {
		filtered := make([]model.Order, 0, len(orders))
		for _, o := range orders {
			if o.Status == status {
				filtered = append(filtered, o)
			}
		}
		orders = filtered
	}

	total := len(orders)

	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return orders[offset:end], total
}

// GetOrder は指定 ID の注文を返す。
func (s *OrderService) GetOrder(id string) (model.Order, bool) {
	return s.repo.FindByID(id)
}
