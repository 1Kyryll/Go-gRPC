package handler

import (
	"net/http"

	"github.com/1kyryll/go-grpc/services/gen/orders"
	"github.com/1kyryll/go-grpc/services/orders/types"
	"github.com/1kyryll/go-grpc/services/util"
)

type OrdersHTTPHandler struct {
	ordersService types.OrderService
}

func NewOrdersHTTPHandler(ordersService types.OrderService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{ordersService: ordersService}
}

func (h *OrdersHTTPHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/order/create", h.CreateOrder)
	router.HandleFunc("/order/get", h.GetOrders)
}

func (h *OrdersHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest
	err := util.ParseJSON(r, &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	order := &orders.Order{
		Id:         1,
		CustomerId: req.GetCustomerId(),
		Items:      req.GetItems(),
	}

	_, err = h.ordersService.CreateOrder(r.Context(), order)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := &orders.CreateOrderResponse{Status: "Success"}
	util.WriteJSON(w, http.StatusOK, res)
}

func (h *OrdersHTTPHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ordersList, err := h.ordersService.GetOrders(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, ordersList)
}
