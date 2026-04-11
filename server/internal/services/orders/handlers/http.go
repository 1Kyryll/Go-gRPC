package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/1kyryll/go-grpc/internal/services/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/services/orders/types"
	"github.com/1kyryll/go-grpc/internal/services/util"
)

type OrdersHTTPHandler struct {
	ordersService types.OrderService
}

func NewOrdersHTTPHandler(ordersService types.OrderService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{ordersService: ordersService}
}

func (h *OrdersHTTPHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/order/create", h.CreateOrder)
	router.HandleFunc("POST /order/{orderId}/done", h.CompleteOrder)
	router.HandleFunc("/ticket/{orderId}", h.GetTickets)
}

func (h *OrdersHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest
	err := util.ParseJSON(r, &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	order := &orders.Order{
		CustomerId: req.GetCustomerId(),
		Items:      req.GetItems(),
	}

	result, err := h.ordersService.CreateOrder(r.Context(), order)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]any{
		"order_id":      result.Order.Id,
		"status":        result.Order.Status,
		"ticket_status": result.TicketStatus,
	})
}

func (h *OrdersHTTPHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order ID: %v", err))
		return
	}

	err = h.ordersService.CompleteOrder(r.Context(), int32(orderID))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (h *OrdersHTTPHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order ID: %v", err))
		return
	}

	tickets, err := h.ordersService.GetTicketsByOrderID(r.Context(), int32(orderID))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if tickets == nil {
		util.WriteJSON(w, http.StatusOK, []struct{}{})
		return
	}
	util.WriteJSON(w, http.StatusOK, tickets)
}
