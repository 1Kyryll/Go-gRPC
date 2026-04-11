package main

import (
	"log"
	"net/http"

	handler "github.com/1kyryll/go-grpc/internal/services/orders/handlers"
	"github.com/1kyryll/go-grpc/internal/services/orders/types"
)

type httpServer struct {
	addr          string
	ordersService types.OrderService
}

func NewHTTPServer(addr string, ordersService types.OrderService) *httpServer {
	return &httpServer{addr: addr, ordersService: ordersService}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	ordersHandler := handler.NewOrdersHTTPHandler(s.ordersService)
	ordersHandler.RegisterRoutes(router)

	log.Printf("HTTP server is running on %s", s.addr)
	return http.ListenAndServe(s.addr, cors(router))
}
