package graph

import (
	"context"
	"fmt"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/gateway/dataloaders"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
)

type orderResolver struct{ *Resolver }
type orderItemResolver struct{ *Resolver }
type customerResolver struct{ *Resolver }
type ticketResolver struct{ *Resolver }

// Order returns OrderResolver implementation.
func (r *Resolver) Order() OrderResolver { return &orderResolver{r} }

// OrderItem returns OrderItemResolver implementation.
func (r *Resolver) OrderItem() OrderItemResolver { return &orderItemResolver{r} }

// Customer returns CustomerResolver implementation.
func (r *Resolver) Customer() CustomerResolver { return &customerResolver{r} }

// Ticket returns TicketResolver implementation.
func (r *Resolver) Ticket() TicketResolver { return &ticketResolver{r} }

// ─── Order fields ───────────────────────────────────────────────────────────

// Customer is the resolver for the customer field.
func (r *orderResolver) Customer(ctx context.Context, obj *model.Order) (*model.Customer, error) {
	c, err := dataloaders.For(ctx).CustomerByID.Load(ctx, obj.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("loading customer %d: %w", obj.CustomerID, err)
	}
	return mapCustomer(c), nil
}

// Items is the resolver for the items field.
func (r *orderResolver) Items(ctx context.Context, obj *model.Order) ([]*model.OrderItem, error) {
	orderID, err := parseID(obj.ID)
	if err != nil {
		return nil, err
	}

	rows, err := dataloaders.For(ctx).OrderItemsByOrderID.Load(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("loading order items for order %s: %w", obj.ID, err)
	}

	items := make([]*model.OrderItem, len(rows))
	for i, row := range rows {
		items[i] = &model.OrderItem{
			ID:                  fmt.Sprintf("%d", row.ID),
			MenuItemID:          row.MenuItemID,
			Quantity:            int(row.Quantity),
			SpecialInstructions: model.TextToStringPtr(row.SpecialInstructions),
		}
	}
	return items, nil
}

// Ticket is the resolver for the ticket field.
func (r *orderResolver) Ticket(ctx context.Context, obj *model.Order) (*model.Ticket, error) {
	orderID, err := parseID(obj.ID)
	if err != nil {
		return nil, err
	}

	t, err := dataloaders.For(ctx).TicketByOrderID.Load(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("loading ticket for order %s: %w", obj.ID, err)
	}
	if t == nil {
		return nil, nil
	}

	return &model.Ticket{
		ID:        fmt.Sprintf("%d", t.ID),
		OrderID:   t.OrderID,
		Status:    model.TicketStatus(t.Status),
		CreatedAt: model.TimestampToTime(t.CreatedAt),
		UpdatedAt: model.TimestampToTime(t.UpdatedAt),
	}, nil
}

// TotalPrice is the resolver for the totalPrice field.
func (r *orderResolver) TotalPrice(ctx context.Context, obj *model.Order) (float64, error) {
	orderID, err := parseID(obj.ID)
	if err != nil {
		return 0, err
	}

	rows, err := dataloaders.For(ctx).OrderItemsByOrderID.Load(ctx, orderID)
	if err != nil {
		return 0, fmt.Errorf("loading items for total price of order %s: %w", obj.ID, err)
	}

	var total float64
	for _, item := range rows {
		mi, err := dataloaders.For(ctx).MenuItemByID.Load(ctx, item.MenuItemID)
		if err != nil {
			return 0, fmt.Errorf("loading menu item %d: %w", item.MenuItemID, err)
		}
		price, _ := mi.Price.Float64Value()
		total += price.Float64 * float64(item.Quantity)
	}
	return total, nil
}

// ─── OrderItem fields ───────────────────────────────────────────────────────

// MenuItem is the resolver for the menuItem field.
func (r *orderItemResolver) MenuItem(ctx context.Context, obj *model.OrderItem) (model.MenuItem, error) {
	mi, err := dataloaders.For(ctx).MenuItemByID.Load(ctx, obj.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("loading menu item %d: %w", obj.MenuItemID, err)
	}
	return model.MapMenuItem(mi), nil
}

// Subtotal is the resolver for the subtotal field.
func (r *orderItemResolver) Subtotal(ctx context.Context, obj *model.OrderItem) (float64, error) {
	mi, err := dataloaders.For(ctx).MenuItemByID.Load(ctx, obj.MenuItemID)
	if err != nil {
		return 0, fmt.Errorf("loading menu item %d for subtotal: %w", obj.MenuItemID, err)
	}
	price, _ := mi.Price.Float64Value()
	return price.Float64 * float64(obj.Quantity), nil
}

// ─── Customer fields ────────────────────────────────────────────────────────

// Orders is the resolver for the orders field.
func (r *customerResolver) Orders(ctx context.Context, obj *model.Customer, first *int, after *string) (*model.OrderConnection, error) {
	customerID, err := parseID(obj.ID)
	if err != nil {
		return nil, err
	}

	limit := int32(20)
	if first != nil {
		limit = int32(*first)
	}

	var afterID int32
	if after != nil {
		afterID, err = decodeCursor(*after)
		if err != nil {
			return nil, err
		}
	}

	rows, err := r.queries.GetOrdersByCustomerIDPaginated(ctx, sqlc.GetOrdersByCustomerIDPaginatedParams{
		CustomerID: customerID,
		AfterID:    afterID,
		PageLimit:  limit + 1,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching customer orders: %w", err)
	}

	hasNext := len(rows) > int(limit)
	if hasNext {
		rows = rows[:limit]
	}

	edges := make([]*model.OrderEdge, len(rows))
	for i, row := range rows {
		edges[i] = &model.OrderEdge{
			Cursor: encodeCursor(row.ID),
			Node:   mapOrder(row),
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	totalCount, err := r.queries.CountOrdersByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("counting customer orders: %w", err)
	}

	return &model.OrderConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNext,
			HasPreviousPage: afterID > 0,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: int(totalCount),
	}, nil
}

// ─── Ticket fields ──────────────────────────────────────────────────────────

// Order is the resolver for the order field.
func (r *ticketResolver) Order(ctx context.Context, obj *model.Ticket) (*model.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, obj.OrderID)
	if err != nil {
		return nil, fmt.Errorf("loading order %d for ticket: %w", obj.OrderID, err)
	}
	return mapOrder(row), nil
}
