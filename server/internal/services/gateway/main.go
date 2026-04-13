package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/gateway/dataloaders"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load()

	// Connect to database
	pool, err := pgxpool.New(context.Background(), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)

	// Connect to gRPC order service
	grpcAddr := os.Getenv("ORDER_SERVICE_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":9000"
	}
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient := orders.NewOrderServiceClient(conn)

	// Create Resolver
	resolver := graph.NewResolver(queries, grpcClient)

	// Create gqlgen server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	// Add transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	mux.Handle("/graphql", dataloaders.Middleware(queries, srv))

	log.Println("GraphQL gateway is running on :8082, playground available at http://localhost:8082/")
	if err := http.ListenAndServe(":8082", cors(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
