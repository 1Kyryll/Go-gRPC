package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/middleware"
	kitchenHandlers "github.com/1kyryll/go-grpc/internal/services/kitchen/handlers"
	kitchenService "github.com/1kyryll/go-grpc/internal/services/kitchen/services"
)

// enrichedOrder is what we broadcast over SSE — the gRPC Order enriched with item names.
type enrichedOrder struct {
	ID     int32    `json:"id"`
	UserID int32    `json:"user_id"`
	Status string   `json:"status"`
	Items  []string `json:"items"`
}

func main() {
	godotenv.Load()

	// Connect to database
	pool, err := pgxpool.New(context.Background(), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	svc := kitchenService.NewKitchenService(queries)

	// Set up Kafka consumer
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(kafkaBrokers, ","),
		Topic:    "orders.events",
		GroupID:  "kitchen-service",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	// Create SSE server
	sse := NewSSEServer()

	// Consume order events from Kafka
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumeOrderEvents(ctx, reader, svc, sse)
	log.Println("Kitchen Kafka consumer started on topic orders.events")

	// Start HTTP server with CORS: SSE + ticket endpoints
	mux := http.NewServeMux()
	mux.Handle("/stream", sse)

	kitchenHTTP := kitchenHandlers.NewKitchenHTTPHandler(svc)
	kitchenHTTP.RegisterRoutes(mux)

	go func() {
		log.Println("Kitchen HTTP server is running on :8081")
		if err := http.ListenAndServe(":8081", middleware.AuthMiddleware()(cors(mux))); err != nil {
			log.Fatalf("Kitchen HTTP server failed: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Kitchen service is running. Press Ctrl+C to exit.")
	<-sig
	cancel()
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
