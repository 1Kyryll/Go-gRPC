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

	"github.com/1kyryll/go-grpc/internal/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/common/gen/user"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/middleware"
	"github.com/1kyryll/go-grpc/internal/services/gateway/dataloaders"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load()

	pool, err := pgxpool.New(context.Background(), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)

	ordersAddr := os.Getenv("ORDERS_GRPC_ADDR")
	if ordersAddr == "" {
		ordersAddr = ":9000"
	}
	ordersConn, err := grpc.NewClient(ordersAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Orders gRPC server: %v", err)
	}
	defer ordersConn.Close()

	userAddr := os.Getenv("USER_GRPC_ADDR")
	if userAddr == "" {
		userAddr = ":9001"
	}
	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User gRPC server: %v", err)
	}
	defer userConn.Close()

	grpcClient := orders.NewOrderServiceClient(ordersConn)
	userGrpcClient := user.NewUserServiceClient(userConn)

	resolver := graph.NewResolver(queries, grpcClient, userGrpcClient)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

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

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	mux.Handle("/graphql", dataloaders.Middleware(queries, srv))

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("GraphQL gateway running on :%s", port)
	log.Printf("Playground: http://localhost:%s/", port)
	if err := http.ListenAndServe(":"+port, middleware.AuthMiddleware()(cors(mux))); err != nil {
		log.Fatalf("Gateway server failed: %v", err)
	}
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
