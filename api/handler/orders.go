package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/order-app/api/model"
	"example.com/order-app/api/service"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

// validStatuses は status クエリパラメータに指定できる値。
var validStatuses = map[string]bool{
	"pending":  true,
	"shipped":  true,
	"canceled": true,
}

// OrderHandler は注文関連の HTTP エンドポイントを提供する。
type OrderHandler struct {
	svc *service.OrderService
}

// NewOrderHandler は OrderHandler を生成する。
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// RegisterRoutes は mux に注文関連のルートを登録する。
func (h *OrderHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/orders", h.ListOrders)
	mux.HandleFunc("GET /api/orders/{id}", h.GetOrder)
}

type ordersResponse struct {
	Orders []model.Order `json:"orders"`
	Total  int           `json:"total"`
}

// ListOrders は GET /api/orders を処理する。
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status != "" && !validStatuses[status] {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid parameter"})
		return
	}

	limit := defaultLimit
	offset := defaultOffset

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}

	orders, total := h.svc.ListOrders(status, limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordersResponse{Orders: orders, Total: total})
}

// GetOrder は GET /api/orders/{id} を処理する。
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	order, ok := h.svc.GetOrder(id)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
