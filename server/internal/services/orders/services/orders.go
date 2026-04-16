package service

import (
	"context"
	"sync"

	"github.com/1kyryll/go-grpc/internal/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrdersService struct {
	mu          sync.Mutex
	queries     *sqlc.Queries
	subscribers map[chan *orders.Order]struct{}
}

func NewOrdersService(queries *sqlc.Queries) *OrdersService {
	return &OrdersService{
		queries:     queries,
		subscribers: make(map[chan *orders.Order]struct{}),
	}
}

func (s *OrdersService) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.Order, error) {
	orderRow, err := s.queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID: req.CustomerId,
		Status: "PENDING",
	})
	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		_, err := s.queries.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID:    orderRow.ID,
			MenuItemID: item.MenuItemId,
			Quantity:   item.Quantity,
			SpecialInstructions: pgtype.Text{
				String: item.SpecialInstructions,
				Valid:  item.SpecialInstructions != "",
			},
		})
		if err != nil {
			return nil, err
		}
	}

	order := &orders.Order{
		Id:         orderRow.ID,
		CustomerId: req.CustomerId,
		Status:     orderRow.Status,
	}

	// Broadcast order to kitchen subscribers
	s.mu.Lock()
	subs := make([]chan *orders.Order, 0, len(s.subscribers))
	for ch := range s.subscribers {
		subs = append(subs, ch)
	}
	s.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- order:
		default:
		}
	}

	return order, nil
}

func (s *OrdersService) GetOrders(ctx context.Context, userID int32) ([]*orders.Order, error) {
	dbOrders, err := s.queries.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*orders.Order, len(dbOrders))
	for i, o := range dbOrders {
		result[i] = &orders.Order{
			Id:         o.ID,
			CustomerId: o.UserID,
			Status:     o.Status,
		}
	}

	return result, nil
}

func (s *OrdersService) Subscribe(ctx context.Context) chan *orders.Order {
	ch := make(chan *orders.Order, 16)
	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()
	return ch
}

func (s *OrdersService) Unsubscribe(ch chan *orders.Order) {
	s.mu.Lock()
	delete(s.subscribers, ch)
	s.mu.Unlock()
	close(ch)
}
