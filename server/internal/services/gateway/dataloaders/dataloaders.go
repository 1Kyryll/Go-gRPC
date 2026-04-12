package dataloaders

import (
	"context"
	"net/http"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
)

type ctxKey string

const loadersKey ctxKey = "dataloaders"

type Loaders struct {
	CustomerByID        *CustomerLoader
	MenuItemByID        *MenuItemLoader
	TicketByOrderID     *TicketLoader
	OrderItemsByOrderID *OrderItemsLoader
}

func NewLoaders(queries *sqlc.Queries) *Loaders {
	return &Loaders{
		CustomerByID:        NewCustomerLoader(queries),
		MenuItemByID:        NewMenuItemLoader(queries),
		TicketByOrderID:     NewTicketLoader(queries),
		OrderItemsByOrderID: NewOrderItemsLoader(queries),
	}
}

func Middleware(queries *sqlc.Queries, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := NewLoaders(queries)
		ctx := context.WithValue(r.Context(), loadersKey, loaders)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
