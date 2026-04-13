package graph

import (
	"context"
	"fmt"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
)

type queryResolver struct{ *Resolver }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Order is the resolver for the order field.
func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	orderID, err := parseID(id)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order %s not found: %w", id, err)
	}

	return mapOrder(row), nil
}

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context, first *int, after *string, status *model.OrderStatus) (*model.OrderConnection, error) {
	limit := int32(20)
	if first != nil {
		limit = int32(*first)
	}

	var afterID int32
	var err error
	if after != nil {
		afterID, err = decodeCursor(*after)
		if err != nil {
			return nil, err
		}
	}

	var statusStr string
	if status != nil {
		statusStr = status.String()
	}

	rows, err := r.queries.GetOrdersPaginated(ctx, sqlc.GetOrdersPaginatedParams{
		AfterID:   afterID,
		Status:    statusStr,
		PageLimit: limit + 1,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching orders: %w", err)
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

	totalCount, err := r.queries.CountOrders(ctx, statusStr)
	if err != nil {
		return nil, fmt.Errorf("counting orders: %w", err)
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

// Customer is the resolver for the customer field.
func (r *queryResolver) Customer(ctx context.Context, id string) (*model.Customer, error) {
	customerID, err := parseID(id)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer %s not found: %w", id, err)
	}

	return mapCustomer(row), nil
}

// MenuItems is the resolver for the menuItems field.
func (r *queryResolver) MenuItems(ctx context.Context, first *int, after *string, category *model.MenuCategory) (*model.MenuItemConnection, error) {
	limit := int32(20)
	if first != nil {
		limit = int32(*first)
	}

	var afterID int32
	var err error
	if after != nil {
		afterID, err = decodeCursor(*after)
		if err != nil {
			return nil, err
		}
	}

	var categoryStr string
	if category != nil {
		categoryStr = category.String()
	}

	rows, err := r.queries.GetMenuItemsPaginated(ctx, sqlc.GetMenuItemsPaginatedParams{
		AfterID:   afterID,
		Category:  categoryStr,
		PageLimit: limit + 1,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching menu items: %w", err)
	}

	hasNext := len(rows) > int(limit)
	if hasNext {
		rows = rows[:limit]
	}

	edges := make([]*model.MenuItemEdge, len(rows))
	for i, row := range rows {
		edges[i] = &model.MenuItemEdge{
			Cursor: encodeCursor(row.ID),
			Node:   model.MapMenuItem(row),
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	totalCount, err := r.queries.CountMenuItems(ctx, categoryStr)
	if err != nil {
		return nil, fmt.Errorf("counting menu items: %w", err)
	}

	return &model.MenuItemConnection{
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

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, query string) ([]model.SearchResult, error) {
	searchTerm := model.StringToText(query)

	var results []model.SearchResult

	menuItems, err := r.queries.SearchMenuItems(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("searching menu items: %w", err)
	}
	for _, item := range menuItems {
		results = append(results, model.MapMenuItem(item).(model.SearchResult))
	}

	customers, err := r.queries.SearchCustomers(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("searching customers: %w", err)
	}
	for _, c := range customers {
		results = append(results, mapCustomer(c))
	}

	searchOrders, err := r.queries.SearchOrders(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("searching orders: %w", err)
	}
	for _, o := range searchOrders {
		results = append(results, mapOrder(o))
	}

	return results, nil
}
