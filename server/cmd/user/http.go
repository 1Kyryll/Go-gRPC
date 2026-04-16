package main

import (
	"log"
	"net/http"

	handler "github.com/1kyryll/go-grpc/internal/services/user/handlers"
	"github.com/1kyryll/go-grpc/internal/services/user/types"
)

type httpServer struct {
	addr        string
	userService types.UserService
}

func NewHTTPServer(addr string, userService types.UserService) *httpServer {
	return &httpServer{addr: addr, userService: userService}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	userHandler := handler.NewUserHTTPHandler(s.userService)
	userHandler.RegisterRoutes(router)

	log.Printf("User HTTP server is running on %s", s.addr)
	return http.ListenAndServe(s.addr, cors(router))
}
