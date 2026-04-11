package types

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/kitchen"
)

type TicketService interface {
	CreateTicket(context.Context, *kitchen.Ticket) (*kitchen.Ticket, error)
	GetTickets(context.Context, int32) ([]*kitchen.Ticket, error)
}
