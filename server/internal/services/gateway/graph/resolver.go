package graph

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"

	"github.com/1kyryll/go-grpc/internal/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/common/gen/user"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
)

type Resolver struct {
	queries        *sqlc.Queries
	grpcClient     orders.OrderServiceClient
	userGrpcClient user.UserServiceClient

	// For subscriptions
	mu                sync.Mutex
	orderSubscribers  map[chan *orders.Order]struct{}
	statusSubscribers map[chan *orders.Order]struct{}
}

func NewResolver(queries *sqlc.Queries, grpcClient orders.OrderServiceClient, userGrpcClient user.UserServiceClient) *Resolver {
	return &Resolver{
		queries:           queries,
		grpcClient:        grpcClient,
		userGrpcClient:    userGrpcClient,
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
		ID:        fmt.Sprintf("%d", o.ID),
		UserID:    o.UserID,
		Status:    model.OrderStatus(o.Status),
		CreatedAt: model.TimestampToTime(o.CreatedAt),
		UpdatedAt: model.TimestampToTime(o.UpdatedAt),
	}
}

func mapUser(u sqlc.User) *model.User {
	return &model.User{
		ID:        fmt.Sprintf("%d", u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Phone:     model.TextToStringPtr(u.Phone),
		CreatedAt: model.TimestampToTime(u.CreatedAt),
		UpdatedAt: model.TimestampToTime(u.UpdatedAt),
	}
}

func parseID(id string) (int32, error) {
	n, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID %q: %w", id, err)
	}
	return int32(n), nil
}

// broadcastNewOrder sends a new order to all orderCreated subscribers.
func (r *Resolver) broadcastNewOrder(order *model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderID, _ := strconv.ParseInt(order.ID, 10, 32)
	protoOrder := &orders.Order{
		Id:         int32(orderID),
		CustomerId: order.UserID,
		Status:     string(order.Status),
	}

	for ch := range r.orderSubscribers {
		select {
		case ch <- protoOrder:
		default:
		}
	}
}

// broadcastStatusChange sends an updated order to all orderStatusChanged subscribers.
func (r *Resolver) broadcastStatusChange(order *model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderID, _ := strconv.ParseInt(order.ID, 10, 32)
	protoOrder := &orders.Order{
		Id:         int32(orderID),
		CustomerId: order.UserID,
		Status:     string(order.Status),
	}

	for ch := range r.statusSubscribers {
		select {
		case ch <- protoOrder:
		default:
		}
	}
}
