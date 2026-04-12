package dataloaders

import (
	"context"
	"time"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/graph-gophers/dataloader/v7"
)

type CustomerLoader struct {
	loader *dataloader.Loader[int32, sqlc.Customer]
}

func NewCustomerLoader(queries *sqlc.Queries) *CustomerLoader {
	batchFn := func(ctx context.Context, keys []int32) []*dataloader.Result[sqlc.Customer] {
		customers, err := queries.GetCustomersByIDs(ctx, keys)

		customerMap := make(map[int32]sqlc.Customer, len(customers))
		if err == nil {
			for _, c := range customers {
				customerMap[c.ID] = c
			}
		}

		results := make([]*dataloader.Result[sqlc.Customer], len(keys))
		for i, key := range keys {
			if customer, ok := customerMap[key]; ok {
				results[i] = &dataloader.Result[sqlc.Customer]{Data: customer, Error: nil}
			} else {
				results[i] = &dataloader.Result[sqlc.Customer]{Data: sqlc.Customer{}, Error: err}
			}
		}
		return results
	}

	loader := dataloader.NewBatchedLoader(batchFn, dataloader.WithWait[int32, sqlc.Customer](2*time.Millisecond))
	return &CustomerLoader{loader: loader}
}

func (l *CustomerLoader) Load(ctx context.Context, id int32) (sqlc.Customer, error) {
	thunk := l.loader.Load(ctx, id)
	return thunk()
}
