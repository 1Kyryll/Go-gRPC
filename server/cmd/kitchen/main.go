package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1kyryll/go-grpc/internal/services/common/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Orders gRPC server: %v", err)
	}
	defer conn.Close()

	client := orders.NewOrderServiceClient(conn)
	log.Println("Kitchen connected to Orders service on :9000")

	stream, err := client.StreamCreatedOrders(context.Background(),
		&orders.StreamCreatedOrdersRequest{})
	if err != nil {
		log.Fatalf("Failed to stream created orders: %v", err)
	}

	// Create SSE server
	sse := NewSSEServer()

	// Receive orders from gRPC stream and broadcast to SSE clients
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
			sse.Broadcast(order)
		}
	}()

	// Start HTTP server for SSE
	mux := http.NewServeMux()
	mux.Handle("/stream", sse)

	go func() {
		log.Println("Kitchen SSE server is running on :8081")
		if err := http.ListenAndServe(":8081", mux); err != nil {
			log.Fatalf("SSE server failed: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Kitchen service is running. Press Ctrl+C to exit.")
	<-sig
}
