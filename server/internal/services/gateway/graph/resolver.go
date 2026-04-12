package graph

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
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

func encodeCursor(id int32) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%d", id)))
}

func decodeCursor(cursor string) (int32, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor: %w", err)
	}
	var id int32
	_, err = fmt.Sscanf(string(b), "cursor:%d", &id)
	if err != nil {
		return 0, fmt.Errorf("malformed cursor: %w", err)
	}
	return id, nil
}

func mapOrder(o sqlc.Order) *model.Order {
	return &model.Order{
		ID:         fmt.Sprintf("%d", o.ID),
		CustomerID: o.CustomerID,
		Status:     model.OrderStatus(o.Status),
		CreatedAt:  model.TimestampToTime(o.CreatedAt),
		UpdatedAt:  model.TimestampToTime(o.UpdatedAt),
	}
}

func mapCustomer(c sqlc.Customer) *model.Customer {
	return &model.Customer{
		ID:        fmt.Sprintf("%d", c.ID),
		Name:      c.Name,
		Email:     c.Email,
		Phone:     model.TextToStringPtr(c.Phone),
		CreatedAt: model.TimestampToTime(c.CreatedAt),
	}
}

func parseID(id string) (int32, error) {
	n, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID %q: %w", id, err)
	}
	return int32(n), nil
}
