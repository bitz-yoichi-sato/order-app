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

func TestListOrders_FilterByStatus(t *testing.T) {
	mux := newTestMux()

	cases := []struct {
		status string
		want   int
	}{
		{"pending", 12},
		{"shipped", 9},
		{"canceled", 3},
	}

	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/api/orders?status="+c.status, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status=%s: expected status 200, got %d", c.status, rec.Code)
		}

		var body struct {
			Orders []map[string]any `json:"orders"`
			Total  int              `json:"total"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("status=%s: failed to decode response: %v", c.status, err)
		}
		if len(body.Orders) != c.want {
			t.Errorf("status=%s: expected %d orders, got %d", c.status, c.want, len(body.Orders))
		}
		if body.Total != c.want {
			t.Errorf("status=%s: expected total %d, got %d", c.status, c.want, body.Total)
		}
		for _, o := range body.Orders {
			if o["status"] != c.status {
				t.Errorf("status=%s: got order with status %v", c.status, o["status"])
			}
		}
	}
}

func TestListOrders_StatusOmitted(t *testing.T) {
	mux := newTestMux()

	// status を渡さない場合と空文字の場合は、どちらも従来どおり全件が対象になる。
	for _, target := range []string{"/api/orders", "/api/orders?status="} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected status 200, got %d", target, rec.Code)
		}

		var body struct {
			Orders []map[string]any `json:"orders"`
			Total  int              `json:"total"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: failed to decode response: %v", target, err)
		}
		if len(body.Orders) != 20 {
			t.Errorf("%s: expected 20 orders with default limit, got %d", target, len(body.Orders))
		}
		if body.Total != 24 {
			t.Errorf("%s: expected total 24, got %d", target, body.Total)
		}
	}
}

func TestListOrders_InvalidStatus(t *testing.T) {
	mux := newTestMux()

	// 未知の値と、大文字を含む値はいずれも不正値として扱う。
	for _, target := range []string{"/api/orders?status=unknown", "/api/orders?status=Pending"} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected status 400, got %d", target, rec.Code)
		}

		var body map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: failed to decode response: %v", target, err)
		}
		if body["error"] != "invalid parameter" {
			t.Errorf("%s: expected error %q, got %q", target, "invalid parameter", body["error"])
		}
		if len(body) != 1 {
			t.Errorf("%s: expected only an error field, got %v", target, body)
		}
	}
}

func TestListOrders_ResponseShapeUnchanged(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodGet, "/api/orders?status=shipped", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body) != 2 {
		t.Errorf("expected exactly 2 top-level fields, got %v", body)
	}
	for _, key := range []string{"orders", "total"} {
		if _, ok := body[key]; !ok {
			t.Errorf("expected top-level field %q", key)
		}
	}

	var orders []map[string]any
	if err := json.Unmarshal(body["orders"], &orders); err != nil {
		t.Fatalf("failed to decode orders: %v", err)
	}
	if len(orders) == 0 {
		t.Fatal("expected at least one order to inspect")
	}
	for _, o := range orders {
		if len(o) != 5 {
			t.Errorf("expected 5 fields per order, got %v", o)
		}
		for _, key := range []string{"id", "customer", "status", "amount", "orderedAt"} {
			if _, ok := o[key]; !ok {
				t.Errorf("expected order field %q, got %v", key, o)
			}
		}
	}

	// 該当0件でも orders は null ではなく空配列であること。
	req = httptest.NewRequest(http.MethodGet, "/api/orders?status=canceled&offset=10", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var emptyPage struct {
		Orders json.RawMessage `json:"orders"`
		Total  int             `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &emptyPage); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if string(emptyPage.Orders) != "[]" {
		t.Errorf("expected empty orders to encode as [], got %s", emptyPage.Orders)
	}
	if emptyPage.Total != 3 {
		t.Errorf("expected total 3, got %d", emptyPage.Total)
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
