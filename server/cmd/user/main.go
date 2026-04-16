package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	handler "github.com/1kyryll/go-grpc/internal/services/user/handlers"
	service "github.com/1kyryll/go-grpc/internal/services/user/services"
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
	userService := service.NewUserService(queries)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":9001")
		if err != nil {
			log.Fatalf("Failed to listen on :9001: %v", err)
		}
		grpcServer := grpc.NewServer()
		handler.NewUserGrpcService(grpcServer, userService)
		log.Println("User gRPC server is running on :9001")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	httpServer := NewHTTPServer(":8083", userService)
	httpServer.Run()
}
