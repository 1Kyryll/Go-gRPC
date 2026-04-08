package main

import (
	"log"
	"net"

	handler "github.com/1kyryll/go-grpc/services/orders/handlers"
	service "github.com/1kyryll/go-grpc/services/orders/services"
	"google.golang.org/grpc"
)

func main() {
	ordersService := service.NewOrdersService()

	// Start gRPC server for Kitchen
	go func() {
		lis, err := net.Listen("tcp", ":9000")
		if err != nil {
			log.Fatalf("Failed to listen on :9000: %v", err)
		}
		grpcServer := grpc.NewServer()
		handler.NewOrdersGrpcService(grpcServer, ordersService)
		log.Println("gRPC server is running on :9000")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server for frontend
	httpServer := NewHTTPServer(":8080", ordersService)
	httpServer.Run()
}
