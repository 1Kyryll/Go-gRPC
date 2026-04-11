package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/1kyryll/go-grpc/internal/services/kitchen/types"
	"github.com/1kyryll/go-grpc/internal/services/util"
)

type KitchenHTTPHandler struct {
	ticketService types.TicketService
}

func NewKitchenHTTPHandler(ticketService types.TicketService) *KitchenHTTPHandler {
	return &KitchenHTTPHandler{ticketService: ticketService}
}

func (h *KitchenHTTPHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/ticket/{orderId}", h.GetTickets)
	router.HandleFunc("POST /order/{orderId}/done", h.CompleteOrder)
}

func (h *KitchenHTTPHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order ID: %v", err))
		return
	}

	tickets, err := h.ticketService.GetTicketsByOrderID(r.Context(), int32(orderID))
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

func (h *KitchenHTTPHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order ID: %v", err))
		return
	}

	err = h.ticketService.CompleteOrder(r.Context(), int32(orderID))
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}
