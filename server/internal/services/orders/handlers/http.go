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
	router.HandleFunc("/order/get/{id}", h.GetOrders)
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
	customerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid customer ID: %v", err))
		return
	}

	ordersList, err := h.ordersService.GetOrders(r.Context(), int32(customerID))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if ordersList == nil {
		ordersList = []*orders.Order{}
	}
	util.WriteJSON(w, http.StatusOK, ordersList)
}
