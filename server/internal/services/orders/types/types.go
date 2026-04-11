package types

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) (*orders.Order, error)
	GetOrders(context.Context, int32) ([]*orders.Order, error)
	Subscribe(context.Context) chan *orders.Order
	Unsubscribe(chan *orders.Order)
}
