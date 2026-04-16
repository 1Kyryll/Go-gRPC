package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/1kyryll/go-grpc/internal/common/gen/user"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	queries *sqlc.Queries
}

func NewUserService(queries *sqlc.Queries) *UserServiceImpl {
	return &UserServiceImpl{queries: queries}
}

func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	_, err := s.queries.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("username %q is already taken", req.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	dbUser, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Phone:        pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		Role:         "CUSTOMER",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, expiredAt, err := genToken(dbUser.ID, dbUser.Username, dbUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &user.RegisterResponse{
		User: &user.User{
			Id:       dbUser.ID,
			Username: dbUser.Username,
			Email:    dbUser.Email,
			Phone:    textToString(dbUser.Phone),
			Role:     dbUser.Role,
		},
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	}, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	dbUser, err := s.queries.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	accessToken, expiredAt, err := genToken(dbUser.ID, dbUser.Username, dbUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &user.LoginResponse{
		User: &user.User{
			Id:       dbUser.ID,
			Username: dbUser.Username,
			Email:    dbUser.Email,
			Phone:    textToString(dbUser.Phone),
			Role:     dbUser.Role,
		},
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	}, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id int32) (*user.User, error) {
	dbUser, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user.User{
		Id:       dbUser.ID,
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Phone:    textToString(dbUser.Phone),
		Role:     dbUser.Role,
	}, nil
}

func genToken(userID int32, username, role string) (string, int64, error) {
	expiredAt := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      fmt.Sprintf("%d", userID),
		"username": username,
		"role":     role,
		"exp":      expiredAt.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "Restaurant",
	})

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}

	return accessToken, expiredAt.Unix(), nil
}

func textToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}
