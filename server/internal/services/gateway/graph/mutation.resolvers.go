package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
)

type mutationResolver struct{ *Resolver }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.CreateOrderInput) (*model.CreateOrderPayload, error) {
	var validationErrors []*model.ValidationError

	if len(input.Items) == 0 {
		validationErrors = append(validationErrors, &model.ValidationError{
			Field:   "items",
			Message: "at least one item is required",
		})
	}

	customerID, err := strconv.Atoi(input.CustomerID)
	if err != nil {
		validationErrors = append(validationErrors, &model.ValidationError{
			Field:   "customerId",
			Message: "invalid customer ID",
		})
	}

	for i, item := range input.Items {
		if item.Quantity <= 0 {
			validationErrors = append(validationErrors, &model.ValidationError{
				Field:   fmt.Sprintf("items[%d].quantity", i),
				Message: "quantity must be greater than 0",
			})
		}

		menuItemID, err := strconv.Atoi(item.MenuItemID)
		if err != nil {
			validationErrors = append(validationErrors, &model.ValidationError{
				Field:   fmt.Sprintf("items[%d].menuItemId", i),
				Message: "invalid menu item ID",
			})
			continue
		}

		dbItem, err := r.queries.GetMenuItemByID(ctx, int32(menuItemID))
		if err != nil {
			validationErrors = append(validationErrors, &model.ValidationError{
				Field:   fmt.Sprintf("items[%d].menuItemId", i),
				Message: fmt.Sprintf("menu item %d not found", menuItemID),
			})
			continue
		}
		if dbItem.IsAvailable.Valid && !dbItem.IsAvailable.Bool {
			validationErrors = append(validationErrors, &model.ValidationError{
				Field:   fmt.Sprintf("items[%d].menuItemId", i),
				Message: fmt.Sprintf("menu item '%s' is not available", dbItem.Name),
			})
		}
	}

	if len(validationErrors) > 0 {
		return &model.CreateOrderPayload{Errors: validationErrors}, nil
	}

	grpcItems := make([]*orders.OrderItemInput, len(input.Items))
	for i, item := range input.Items {
		menuItemID, _ := strconv.Atoi(item.MenuItemID)
		specialInstructions := ""
		if item.SpecialInstructions != nil {
			specialInstructions = *item.SpecialInstructions
		}
		grpcItems[i] = &orders.OrderItemInput{
			MenuItemId:          int32(menuItemID),
			Quantity:            int32(item.Quantity),
			SpecialInstructions: specialInstructions,
		}
	}

	_, err = r.grpcClient.CreateOrder(ctx, &orders.CreateOrderRequest{
		CustomerId: int32(customerID),
		Items:      grpcItems,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order via gRPC: %w", err)
	}

	dbOrders, err := r.queries.GetOrders(ctx, int32(customerID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created order: %w", err)
	}
	if len(dbOrders) == 0 {
		return nil, fmt.Errorf("order created but not found")
	}

	latest := dbOrders[len(dbOrders)-1]
	order := mapOrder(latest)

	r.broadcastNewOrder(order)

	return &model.CreateOrderPayload{Order: order}, nil
}

// UpdateOrderStatus is the resolver for the updateOrderStatus field.
func (r *mutationResolver) UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) (*model.Order, error) {
	orderID, err := parseID(id)
	if err != nil {
		return nil, err
	}

	dbOrder, err := r.queries.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     orderID,
		Status: status.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	order := mapOrder(dbOrder)
	r.broadcastStatusChange(order)
	return order, nil
}

// CompleteOrder is the resolver for the completeOrder field.
func (r *mutationResolver) CompleteOrder(ctx context.Context, orderID string) (*model.Order, error) {
	id, err := parseID(orderID)
	if err != nil {
		return nil, err
	}

	if err := r.queries.CompleteTicketByOrderID(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to complete ticket: %w", err)
	}
	if err := r.queries.CompleteOrder(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to complete order: %w", err)
	}

	dbOrder, err := r.queries.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch completed order: %w", err)
	}

	order := mapOrder(dbOrder)
	r.broadcastStatusChange(order)
	return order, nil
}

// CancelOrder is the resolver for the cancelOrder field.
func (r *mutationResolver) CancelOrder(ctx context.Context, orderID string) (*model.Order, error) {
	id, err := parseID(orderID)
	if err != nil {
		return nil, err
	}

	dbOrder, err := r.queries.CancelOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	order := mapOrder(dbOrder)
	r.broadcastStatusChange(order)
	return order, nil
}
