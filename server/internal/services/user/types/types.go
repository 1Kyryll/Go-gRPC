package types

import (
	"context"

	"github.com/1kyryll/go-grpc/internal/common/gen/user"
)

type UserService interface {
	Register(context.Context, *user.RegisterRequest) (*user.RegisterResponse, error)
	Login(context.Context, *user.LoginRequest) (*user.LoginResponse, error)
	GetUser(context.Context, int32) (*user.User, error)
}
