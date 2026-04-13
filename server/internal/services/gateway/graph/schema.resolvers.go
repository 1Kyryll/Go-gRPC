package graph

import (
	"context"
	"fmt"

	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
)

type subscriptionResolver struct{ *Resolver }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// OrderCreated is the resolver for the orderCreated field.
func (r *subscriptionResolver) OrderCreated(ctx context.Context) (<-chan *model.Order, error) {
	panic(fmt.Errorf("not implemented: OrderCreated - orderCreated"))
}

// OrderStatusChanged is the resolver for the orderStatusChanged field.
func (r *subscriptionResolver) OrderStatusChanged(ctx context.Context, orderID *string) (<-chan *model.Order, error) {
	panic(fmt.Errorf("not implemented: OrderStatusChanged - orderStatusChanged"))
}
