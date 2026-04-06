package service

import (
	"context"

	"github.com/1kyryll/go-grpc/services/gen/orders"
)

var ordersDb = make([]*orders.Order, 0)

type OrdersService struct{}

func NewOrdersService() *OrdersService {
	return &OrdersService{}
}

func (s *OrdersService) CreateOrder(ctx context.Context, order *orders.Order) (*orders.Order, error) {
	ordersDb = append(ordersDb, order)
	return order, nil
}

func (s *OrdersService) GetOrders(ctx context.Context) ([]*orders.Order, error) {
	return ordersDb, nil
}
