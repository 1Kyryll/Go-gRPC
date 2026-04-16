package middleware

import (
	"context"
	"fmt"
)

const (
	RoleCustomer     = "CUSTOMER"
	RoleKitchenStaff = "KITCHEN_STAFF"
)

func RequireAuth(ctx context.Context) (*AuthUser, error) {
	user, err := GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, fmt.Errorf("Authentication required")
	}

	return user, nil
}

func RequireRole(ctx context.Context, role string) (*AuthUser, error) {
	user, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if user.Role != role {
		return nil, fmt.Errorf("%s role required", role)
	}

	return user, nil
}
