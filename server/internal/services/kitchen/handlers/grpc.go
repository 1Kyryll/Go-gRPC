package handlers

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/kitchen"
	"github.com/1kyryll/go-grpc/internal/services/kitchen/types"
	"google.golang.org/grpc"
)

type KitchenGrpcHandler struct {
	ticketService types.TicketService
	kitchen.UnimplementedTicketServiceServer
}

func NewKitchenGrpcService(grpc *grpc.Server, ticketService types.TicketService) {
	gRPCHandler := &KitchenGrpcHandler{ticketService: ticketService}

	//register the TicketServiceServer
	kitchen.RegisterTicketServiceServer(grpc, gRPCHandler)
}

func (h *KitchenGrpcHandler) CreateTicket(ctx context.Context, req *kitchen.CreateTicketRequest) (*kitchen.CreateTicketResponse, error) {
	ticket := &kitchen.Ticket{
		OrderId: req.OrderId,
		Status:  "Cooking",
	}
	_, err := h.ticketService.CreateTicket(ctx, ticket)
	if err != nil {
		return nil, err
	}

	return &kitchen.CreateTicketResponse{Status: "Success"}, nil
}

func (h *KitchenGrpcHandler) GetTicketsByOrderID(ctx context.Context, req *kitchen.GetTicketsByOrderIDRequest) (*kitchen.GetTicketsByOrderIDResponse, error) {
	tickets, err := h.ticketService.GetTickets(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &kitchen.GetTicketsByOrderIDResponse{Tickets: tickets}, nil
}
