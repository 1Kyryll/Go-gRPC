package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/kitchen"
	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	kitchenHandlers "github.com/1kyryll/go-grpc/internal/services/kitchen/handlers"
	kitchenService "github.com/1kyryll/go-grpc/internal/services/kitchen/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	// Connect to Orders gRPC server
	grpcAddr := os.Getenv("ORDERS_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":9000"
	}
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Orders gRPC server: %v", err)
	}
	defer conn.Close()

	client := orders.NewOrderServiceClient(conn)
	log.Printf("Kitchen connecting to Orders service on %s", grpcAddr)

	var stream orders.OrderService_StreamCreatedOrdersClient
	for {
		stream, err = client.StreamCreatedOrders(context.Background(),
			&orders.StreamCreatedOrdersRequest{})
		if err == nil {
			break
		}
		log.Printf("Waiting for Orders gRPC server: %v", err)
		time.Sleep(2 * time.Second)
	}
	log.Println("Kitchen connected to Orders stream")

	// Create SSE server
	sse := NewSSEServer()

	// Receive orders from gRPC stream, create ticket, broadcast to SSE
	go func() {
		for {
			order, err := stream.Recv()
			if err == io.EOF {
				log.Println("Stream closed by server")
				return
			}
			if err != nil {
				log.Printf("Error receiving order: %v", err)
				return
			}

			log.Printf("NEW ORDER: id=%d customer=%d items=%v",
				order.Id, order.CustomerId, order.Items)

			// Create a ticket for this order
			ticket := &kitchen.Ticket{
				OrderId: order.Id,
				Status:  "open",
			}
			if _, err := svc.CreateTicket(context.Background(), ticket); err != nil {
				log.Printf("Failed to create ticket for order %d: %v", order.Id, err)
			}

			sse.Broadcast(order)
		}
	}()

	// Start HTTP server with CORS: SSE + ticket endpoints
	mux := http.NewServeMux()
	mux.Handle("/stream", sse)

	kitchenHTTP := kitchenHandlers.NewKitchenHTTPHandler(svc)
	kitchenHTTP.RegisterRoutes(mux)

	go func() {
		log.Println("Kitchen HTTP server is running on :8081")
		if err := http.ListenAndServe(":8081", cors(mux)); err != nil {
			log.Fatalf("Kitchen HTTP server failed: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Kitchen service is running. Press Ctrl+C to exit.")
	<-sig
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
