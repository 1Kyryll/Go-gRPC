package service

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/orders/types"
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

func (s *OrdersService) CreateOrder(ctx context.Context, order *orders.Order) (*types.CreateOrderResult, error) {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return nil, err
	}

	orderRow, err := s.queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		CustomerID: order.CustomerId,
		Items:      itemsJSON,
		Status:     "pending",
	})
	if err != nil {
		return nil, err
	}

	order.Id = orderRow.ID
	order.Status = orderRow.Status

	// Create a ticket for the kitchen
	ticketRow, err := s.queries.CreateTicket(ctx, sqlc.CreateTicketParams{
		OrderID: order.Id,
		Status:  "open",
	})
	if err != nil {
		return nil, err
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

	return &types.CreateOrderResult{
		Order:        order,
		TicketStatus: ticketRow.Status,
	}, nil
}

func (s *OrdersService) GetOrders(ctx context.Context, customerID int32) ([]*orders.Order, error) {
	dbOrders, err := s.queries.GetOrders(ctx, customerID)
	if err != nil {
		return nil, err
	}

	result := make([]*orders.Order, len(dbOrders))
	for i, o := range dbOrders {
		var items []string
		json.Unmarshal(o.Items, &items)
		result[i] = &orders.Order{
			Id:         o.ID,
			CustomerId: o.CustomerID,
			Items:      items,
			Status:     o.Status,
		}
	}

	return result, nil
}

func (s *OrdersService) GetTicketsByOrderID(ctx context.Context, orderID int32) ([]sqlc.Ticket, error) {
	return s.queries.GetTicketsByOrderID(ctx, orderID)
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
