package types

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/kitchen"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
)

type TicketService interface {
	CreateTicket(context.Context, *kitchen.Ticket) (*kitchen.Ticket, error)
	GetTickets(context.Context, int32) ([]*kitchen.Ticket, error)
	CompleteOrder(context.Context, int32) error
	GetTicketsByOrderID(context.Context, int32) ([]sqlc.Ticket, error)
}
