package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	handler "github.com/1kyryll/go-grpc/internal/services/orders/handlers"
	service "github.com/1kyryll/go-grpc/internal/services/orders/services"
	outbox "github.com/1kyryll/go-grpc/internal/services/outbox"
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
	ordersService := service.NewOrdersService(pool, queries)

	// Start gRPC server for Gateway subscriptions
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

	// Start outbox relay (publishes to Kafka)
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "loclahost:9092"
	}
	relay := outbox.NewOutboxRelay(queries, strings.Split(kafkaBrokers, ","))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go relay.Run(ctx)
	log.Println("Outbox relay started")

	// Start HTTP server for frontend
	go func() {
		httpServer := NewHTTPServer(":8080", ordersService)
		httpServer.Run()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Orders service is running. Press Ctrl+C to exit.")
	<-sig
	cancel()
}
