package types

import (
	"context"

	"github.com/1kyryll/go-grpc/services/gen/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) (*orders.Order, error)
	GetOrders(context.Context) ([]*orders.Order, error)
}
