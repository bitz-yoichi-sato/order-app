package main

import (
	"log"
	"net/http"

	"example.com/order-app/api/handler"
	"example.com/order-app/api/repository"
	"example.com/order-app/api/service"
)

func main() {
	repo := repository.NewOrderRepository()
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
