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

// ListOrders は limit/offset に基づき注文一覧と総件数を返す。
func (s *OrderService) ListOrders(limit, offset int) ([]model.Order, int) {
	orders := s.repo.List()
	return orders[offset : offset+limit], len(orders)
}

// GetOrder は指定 ID の注文を返す。
func (s *OrderService) GetOrder(id string) (model.Order, bool) {
	return s.repo.FindByID(id)
}
