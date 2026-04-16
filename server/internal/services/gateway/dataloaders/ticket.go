package dataloaders

import (
	"context"
	"time"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/graph-gophers/dataloader/v7"
)

type TicketLoader struct {
	loader *dataloader.Loader[int32, *sqlc.Ticket]
}

func NewTicketLoader(queries *sqlc.Queries) *TicketLoader {
	batchFn := func(ctx context.Context, keys []int32) []*dataloader.Result[*sqlc.Ticket] {
		tickets, err := queries.GetTicketsByOrderIDs(ctx, keys)

		ticketMap := make(map[int32]sqlc.Ticket, len(tickets))
		if err == nil {
			for _, t := range tickets {
				ticketMap[t.OrderID] = t
			}
		}

		results := make([]*dataloader.Result[*sqlc.Ticket], len(keys))
		for i, key := range keys {
			if t, ok := ticketMap[key]; ok {
				results[i] = &dataloader.Result[*sqlc.Ticket]{Data: &t}
			} else if err != nil {
				results[i] = &dataloader.Result[*sqlc.Ticket]{Error: err}
			} else {
				// No ticket for this order is valid (not an error)
				results[i] = &dataloader.Result[*sqlc.Ticket]{Data: nil}
			}
		}
		return results
	}

	loader := dataloader.NewBatchedLoader(batchFn, dataloader.WithWait[int32, *sqlc.Ticket](2*time.Millisecond))
	return &TicketLoader{loader: loader}
}

func (l *TicketLoader) Load(ctx context.Context, orderID int32) (*sqlc.Ticket, error) {
	thunk := l.loader.Load(ctx, orderID)
	return thunk()
}
