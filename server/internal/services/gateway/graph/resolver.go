package graph

import (
	"sync"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
)

type Resolver struct {
	queries    *sqlc.Queries
	grpcClient orders.OrderServiceClient

	// For subscriptions
	mu                sync.Mutex
	orderSubscribers  map[chan *orders.Order]struct{}
	statusSubscribers map[chan *orders.Order]struct{}
}

func NewResolver(queries *sqlc.Queries, grpcClient orders.OrderServiceClient) *Resolver {
	return &Resolver{
		queries:           queries,
		grpcClient:        grpcClient,
		orderSubscribers:  make(map[chan *orders.Order]struct{}),
		statusSubscribers: make(map[chan *orders.Order]struct{}),
	}
}
