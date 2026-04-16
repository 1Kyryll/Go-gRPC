package dataloaders

import (
	"context"
	"time"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/graph-gophers/dataloader/v7"
)

type UserLoader struct {
	loader *dataloader.Loader[int32, sqlc.User]
}

func NewUserLoader(queries *sqlc.Queries) *UserLoader {
	batchFn := func(ctx context.Context, keys []int32) []*dataloader.Result[sqlc.User] {
		users, err := queries.GetUsersByIDs(ctx, keys)

		userMap := make(map[int32]sqlc.User, len(users))
		if err == nil {
			for _, u := range users {
				userMap[u.ID] = u
			}
		}

		results := make([]*dataloader.Result[sqlc.User], len(keys))
		for i, key := range keys {
			if user, ok := userMap[key]; ok {
				results[i] = &dataloader.Result[sqlc.User]{Data: user, Error: nil}
			} else {
				results[i] = &dataloader.Result[sqlc.User]{Data: sqlc.User{}, Error: err}
			}
		}
		return results
	}

	loader := dataloader.NewBatchedLoader(batchFn, dataloader.WithWait[int32, sqlc.User](2*time.Millisecond))
	return &UserLoader{loader: loader}
}

func (l *UserLoader) Load(ctx context.Context, id int32) (sqlc.User, error) {
	thunk := l.loader.Load(ctx, id)
	return thunk()
}
