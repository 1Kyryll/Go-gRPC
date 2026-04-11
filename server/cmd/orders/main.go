package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	handler "github.com/1kyryll/go-grpc/internal/services/orders/handlers"
	service "github.com/1kyryll/go-grpc/internal/services/orders/services"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	pool, err := pgxpool.New(context.Background(), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	ordersService := service.NewOrdersService(queries)

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
