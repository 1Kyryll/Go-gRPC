package service

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/common/gen/kitchen"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
)

type KitchenService struct {
	queries *sqlc.Queries
}

func NewKitchenService(queries *sqlc.Queries) *KitchenService {
	return &KitchenService{queries: queries}
}

func (s *KitchenService) CreateTicket(ctx context.Context, ticket *kitchen.Ticket) (*kitchen.Ticket, error) {
	row, err := s.queries.CreateTicket(ctx, sqlc.CreateTicketParams{
		OrderID: ticket.OrderId,
		Status:  ticket.Status,
	})
	if err != nil {
		return nil, err
	}

	ticket.Status = row.Status
	return ticket, nil
}

func (s *KitchenService) GetTickets(ctx context.Context, orderID int32) ([]*kitchen.Ticket, error) {
	dbTickets, err := s.queries.GetTicketsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	result := make([]*kitchen.Ticket, len(dbTickets))
	for i, t := range dbTickets {
		result[i] = &kitchen.Ticket{
			Id:      t.ID,
			OrderId: t.OrderID,
			Status:  t.Status,
		}
	}

	return result, nil
}

func (s *KitchenService) CompleteOrder(ctx context.Context, orderID int32) error {
	if err := s.queries.CompleteTicketByOrderID(ctx, orderID); err != nil {
		return err
	}
	return s.queries.CompleteOrder(ctx, orderID)
}

func (s *KitchenService) GetTicketsByOrderID(ctx context.Context, orderID int32) ([]sqlc.Ticket, error) {
	return s.queries.GetTicketsByOrderID(ctx, orderID)
}
