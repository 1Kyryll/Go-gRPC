package handler

import (
	"net/http"

	"github.com/1kyryll/go-grpc/internal/common/gen/user"
	"github.com/1kyryll/go-grpc/internal/services/user/types"
	"github.com/1kyryll/go-grpc/internal/util"
)

type UserHTTPHandler struct {
	userService types.UserService
}

func NewUserHTTPHandler(userService types.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{userService: userService}
}

func (h *UserHTTPHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/auth/register", h.Register)
	router.HandleFunc("/auth/login", h.Login)
}

func (h *UserHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := util.ParseJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, resp)
}

func (h *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := util.ParseJSON(r, &req); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}
