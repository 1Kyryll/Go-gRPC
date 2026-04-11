package service

import (
	"context"
	"sync"

	"github.com/1kyryll/go-grpc/internal/services/common/orders"
)

var ordersDb = make([]*orders.Order, 0)

type OrdersService struct {
	mu          sync.Mutex
	subscribers map[chan *orders.Order]struct{}
}

func NewOrdersService() *OrdersService {
	return &OrdersService{
		subscribers: make(map[chan *orders.Order]struct{}),
	}
}

func (s *OrdersService) CreateOrder(ctx context.Context, order *orders.Order) (*orders.Order, error) {
	s.mu.Lock()
	ordersDb = append(ordersDb, order)

	subs := make([]chan *orders.Order, 0, len(s.subscribers))
	for ch := range s.subscribers {
		subs = append(subs, ch)
	}
	s.mu.Unlock()

	//broadcast to subscribers
	for _, ch := range subs {
		select {
		case ch <- order:
		default:
		}
	}

	return order, nil
}

func (s *OrdersService) GetOrders(ctx context.Context) ([]*orders.Order, error) {
	return ordersDb, nil
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
