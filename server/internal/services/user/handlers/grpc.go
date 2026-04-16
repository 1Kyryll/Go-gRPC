package handler

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/common/gen/user"
	"github.com/1kyryll/go-grpc/internal/services/user/types"
	"google.golang.org/grpc"
)

type UserGrpcHandler struct {
	userService types.UserService
	user.UnimplementedUserServiceServer
}

func NewUserGrpcService(grpc *grpc.Server, userService types.UserService) {
	handler := &UserGrpcHandler{userService: userService}
	user.RegisterUserServiceServer(grpc, handler)
}

func (h *UserGrpcHandler) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	return h.userService.Register(ctx, req)
}

func (h *UserGrpcHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	return h.userService.Login(ctx, req)
}

func (h *UserGrpcHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	return h.userService.GetUser(ctx, req.Id)
}
