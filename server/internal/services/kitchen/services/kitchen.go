package service

import (
	"context"
	"fmt"

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
	err := s.queries.CreateTicketsIdempotent(ctx, sqlc.CreateTicketsIdempotentParams{
		OrderID: ticket.OrderId,
		Status:  ticket.Status,
	})
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *KitchenService) GetEnrichedOrderItems(ctx context.Context, orderID int32) ([]string, error) {
	orderItems, err := s.queries.GetOrderItemsByOrderIDs(ctx, []int32{orderID})
	if err != nil {
		return nil, err
	}

	menuIDs := make([]int32, len(orderItems))
	for i, oi := range orderItems {
		menuIDs[i] = oi.MenuItemID
	}

	menuItems, err := s.queries.GetMenuItemsByIDs(ctx, menuIDs)
	if err != nil {
		return nil, err
	}

	nameMap := make(map[int32]string, len(menuItems))
	for _, mi := range menuItems {
		nameMap[mi.ID] = mi.Name
	}

	var items []string
	for _, oi := range orderItems {
		name := nameMap[oi.MenuItemID]
		if name == "" {
			name = fmt.Sprintf("Item #%d", oi.MenuItemID)
		}
		if oi.Quantity > 1 {
			items = append(items, fmt.Sprintf("%s x%d", name, oi.Quantity))
		} else {
			items = append(items, name)
		}
	}

	return items, nil
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
