package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/order-app/api/repository"
	"example.com/order-app/api/service"
)

func newTestMux() *http.ServeMux {
	repo := repository.NewOrderRepository()
	svc := service.NewOrderService(repo)
	h := NewOrderHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestListOrders_DefaultParams(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body struct {
		Orders []map[string]any `json:"orders"`
		Total  int              `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body.Orders) != 20 {
		t.Errorf("expected 20 orders with default limit, got %d", len(body.Orders))
	}
}

func TestListOrders_Total(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var body struct {
		Orders []map[string]any `json:"orders"`
		Total  int              `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.Total != 24 {
		t.Errorf("expected total 24, got %d", body.Total)
	}
}

func TestGetOrder_Found(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/orders/ord-001", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/orders/does-not-exist", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["error"] == "" {
		t.Errorf("expected error message in response body")
	}
}
