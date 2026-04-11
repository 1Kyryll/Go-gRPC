package handler

import (
	"net/http"

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
		"order_id": result.Id,
		"status":   result.Status,
	})
}
