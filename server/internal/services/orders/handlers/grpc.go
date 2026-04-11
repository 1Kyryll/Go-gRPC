package handler

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/services/common/orders"
	"github.com/1kyryll/go-grpc/internal/services/orders/types"
	"google.golang.org/grpc"
)

type OrdersGrpcHandler struct {
	ordersService types.OrderService
	orders.UnimplementedOrderServiceServer
}

func NewOrdersGrpcService(grpc *grpc.Server, ordersService types.OrderService) {
	gRPCHandler := &OrdersGrpcHandler{ordersService: ordersService}

	//register the OrderServiceServer
	orders.RegisterOrderServiceServer(grpc, gRPCHandler)
}

func (h *OrdersGrpcHandler) StreamCreatedOrders(
	req *orders.StreamCreatedOrdersRequest,
	stream orders.OrderService_StreamCreatedOrdersServer) error {
	ch := h.ordersService.Subscribe(stream.Context())
	defer h.ordersService.Unsubscribe(ch)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case order := <-ch:
			if err := stream.Send(order); err != nil {
				return err
			}
		}
	}
}

func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		Id:         1,
		CustomerId: req.CustomerId,
		Items:      req.Items,
	}

	_, err := h.ordersService.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &orders.CreateOrderResponse{Status: "Success"}, nil
}

func (h *OrdersGrpcHandler) GetOrders(ctx context.Context, req *orders.GetOrdersRequest) (*orders.GetOrdersResponse, error) {
	ordersList, err := h.ordersService.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	return &orders.GetOrdersResponse{Orders: ordersList}, nil
}
