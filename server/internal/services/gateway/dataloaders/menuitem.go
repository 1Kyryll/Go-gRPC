package dataloaders

import (
	"context"
	"fmt"
	"time"

	"github.com/1kyryll/go-grpc/internal/services/common/sqlc"
	"github.com/graph-gophers/dataloader/v7"
)

type MenuItemLoader struct {
	loader *dataloader.Loader[int32, sqlc.MenuItem]
}

func NewMenuItemLoader(queries *sqlc.Queries) *MenuItemLoader {
	batchFn := func(ctx context.Context, keys []int32) []*dataloader.Result[sqlc.MenuItem] {
		items, err := queries.GetMenuItemsByIDs(ctx, keys)

		itemMap := make(map[int32]sqlc.MenuItem, len(items))
		if err == nil {
			for _, item := range items {
				itemMap[item.ID] = item
			}
		}

		results := make([]*dataloader.Result[sqlc.MenuItem], len(keys))
		for i, key := range keys {
			if item, ok := itemMap[key]; ok {
				results[i] = &dataloader.Result[sqlc.MenuItem]{Data: item}
			} else if err != nil {
				results[i] = &dataloader.Result[sqlc.MenuItem]{Error: err}
			} else {
				results[i] = &dataloader.Result[sqlc.MenuItem]{Error: fmt.Errorf("menu item %d not found", key)}
			}
		}
		return results
	}

	loader := dataloader.NewBatchedLoader(batchFn, dataloader.WithWait[int32, sqlc.MenuItem](2*time.Millisecond))
	return &MenuItemLoader{loader: loader}
}

func (l *MenuItemLoader) Load(ctx context.Context, id int32) (sqlc.MenuItem, error) {
	thunk := l.loader.Load(ctx, id)
	return thunk()
}
