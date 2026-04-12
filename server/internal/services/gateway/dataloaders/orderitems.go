package dataloaders

import (
	"context"
	"time"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/graph-gophers/dataloader/v7"
)

type OrderItemsLoader struct {
	loader *dataloader.Loader[int32, []sqlc.OrderItem]
}

func NewOrderItemsLoader(queries *sqlc.Queries) *OrderItemsLoader {
	batchFn := func(ctx context.Context, keys []int32) []*dataloader.Result[[]sqlc.OrderItem] {
		allItems, err := queries.GetOrderItemsByOrderIDs(ctx, keys)

		itemsByOrder := make(map[int32][]sqlc.OrderItem)
		if err == nil {
			for _, item := range allItems {
				itemsByOrder[item.OrderID] = append(itemsByOrder[item.OrderID], item)
			}
		}

		results := make([]*dataloader.Result[[]sqlc.OrderItem], len(keys))
		for i, key := range keys {
			if err != nil {
				results[i] = &dataloader.Result[[]sqlc.OrderItem]{Error: err}
			} else {
				items := itemsByOrder[key]
				if items == nil {
					items = []sqlc.OrderItem{}
				}
				results[i] = &dataloader.Result[[]sqlc.OrderItem]{Data: items}
			}
		}
		return results
	}

	loader := dataloader.NewBatchedLoader(batchFn, dataloader.WithWait[int32, []sqlc.OrderItem](2*time.Millisecond))
	return &OrderItemsLoader{loader: loader}
}

func (l *OrderItemsLoader) Load(ctx context.Context, orderID int32) ([]sqlc.OrderItem, error) {
	thunk := l.loader.Load(ctx, orderID)
	return thunk()
}
