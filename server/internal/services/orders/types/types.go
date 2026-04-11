package types

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
)

type CreateOrderResult struct {
	Order        *orders.Order
	TicketStatus string
}

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) (*CreateOrderResult, error)
	GetOrders(context.Context, int32) ([]*orders.Order, error)
	GetTicketsByOrderID(context.Context, int32) ([]sqlc.Ticket, error)
	Subscribe(context.Context) chan *orders.Order
	Unsubscribe(chan *orders.Order)
}
