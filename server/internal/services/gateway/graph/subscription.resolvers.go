package graph

import (
	"context"
	"io"
	"log"
	"strconv"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/gateway/graph/model"
)

type subscriptionResolver struct{ *Resolver }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// OrderCreated streams newly created orders via the gRPC StreamCreatedOrders RPC.
func (r *subscriptionResolver) OrderCreated(ctx context.Context) (<-chan *model.Order, error) {
	stream, err := r.grpcClient.StreamCreatedOrders(ctx, &orders.StreamCreatedOrdersRequest{})
	if err != nil {
		return nil, err
	}

	ch := make(chan *model.Order, 16)

	go func() {
		defer close(ch)
		for {
			order, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("orderCreated stream error: %v", err)
				return
			}

			gqlOrder := &model.Order{
				ID:         strconv.Itoa(int(order.Id)),
				CustomerID: order.CustomerId,
				Status:     model.OrderStatus(order.Status),
			}

			select {
			case ch <- gqlOrder:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

// OrderStatusChanged streams order status updates, optionally filtered by order ID.
func (r *subscriptionResolver) OrderStatusChanged(ctx context.Context, orderID *string) (<-chan *model.Order, error) {
	ch := make(chan *orders.Order, 16)

	r.mu.Lock()
	r.statusSubscribers[ch] = struct{}{}
	r.mu.Unlock()

	out := make(chan *model.Order, 16)

	go func() {
		defer close(out)
		defer func() {
			r.mu.Lock()
			delete(r.statusSubscribers, ch)
			r.mu.Unlock()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case protoOrder, ok := <-ch:
				if !ok {
					return
				}

				// Filter by orderID if provided
				if orderID != nil {
					targetID, err := strconv.Atoi(*orderID)
					if err != nil {
						continue
					}
					if protoOrder.Id != int32(targetID) {
						continue
					}
				}

				gqlOrder := &model.Order{
					ID:         strconv.Itoa(int(protoOrder.Id)),
					CustomerID: protoOrder.CustomerId,
					Status:     model.OrderStatus(protoOrder.Status),
				}

				select {
				case out <- gqlOrder:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out, nil
}
