package service

import (
	"context"
	"sync"

	"github.com/1kyryll/go-grpc/internal/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
)

type OrdersService struct {
	mu          sync.Mutex
	pool        *pgxpool.Pool
	queries     *sqlc.Queries
	subscribers map[chan *orders.Order]struct{}
}

func NewOrdersService(pool *pgxpool.Pool, queries *sqlc.Queries) *OrdersService {
	return &OrdersService{
		pool:        pool,
		queries:     queries,
		subscribers: make(map[chan *orders.Order]struct{}),
	}
}

func (s *OrdersService) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.Order, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	txQueries := sqlc.New(tx)

	orderRow, err := txQueries.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID: req.CustomerId,
		Status: "PENDING",
	})
	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		_, err := txQueries.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID:    orderRow.ID,
			MenuItemID: item.MenuItemId,
			Quantity:   item.Quantity,
			SpecialInstructions: pgtype.Text{
				String: item.SpecialInstructions,
				Valid:  item.SpecialInstructions != "",
			},
		})
		if err != nil {
			return nil, err
		}
	}

	// Write outbox event in the same transaction
	event := &orders.OrderEvent{
		EventType:  "ORDER_CREATED",
		OrderId:    orderRow.ID,
		CustomerId: req.CustomerId,
		Status:     orderRow.Status,
		Items:      req.Items,
	}
	payload, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}

	err = txQueries.InsertOutboxEvent(ctx, sqlc.InsertOutboxEventParams{
		AggregateID: orderRow.ID,
		EventType:   "ORDER_CREATED",
		Payload:     payload,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	order := &orders.Order{
		Id:         orderRow.ID,
		CustomerId: req.CustomerId,
		Status:     orderRow.Status,
	}

	// Broadcast order to kitchen subscribers
	s.mu.Lock()
	subs := make([]chan *orders.Order, 0, len(s.subscribers))
	for ch := range s.subscribers {
		subs = append(subs, ch)
	}
	s.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- order:
		default:
		}
	}

	return order, nil
}

func (s *OrdersService) GetOrders(ctx context.Context, userID int32) ([]*orders.Order, error) {
	dbOrders, err := s.queries.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*orders.Order, len(dbOrders))
	for i, o := range dbOrders {
		result[i] = &orders.Order{
			Id:         o.ID,
			CustomerId: o.UserID,
			Status:     o.Status,
		}
	}

	return result, nil
}

func (s *OrdersService) Subscribe(ctx context.Context) chan *orders.Order {
	ch := make(chan *orders.Order, 16)
	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()
	return ch
}

func (s *OrdersService) Unsubscribe(ch chan *orders.Order) {
	s.mu.Lock()
	delete(s.subscribers, ch)
	s.mu.Unlock()
	close(ch)
}
