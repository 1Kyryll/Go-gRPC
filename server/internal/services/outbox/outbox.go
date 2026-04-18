package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/segmentio/kafka-go"
)

type OutboxRelay struct {
	queries *sqlc.Queries
	writer  *kafka.Writer
}

func NewOutboxRelay(queries *sqlc.Queries, brokers []string) *OutboxRelay {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  "orders.events",
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.Hash{},
	}
	return &OutboxRelay{
		queries: queries,
		writer:  writer,
	}
}

func (r *OutboxRelay) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer r.writer.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.publishPending(ctx)
		}
	}
}

func (r *OutboxRelay) publishPending(ctx context.Context) {
	events, err := r.queries.GetUnpublishedOutboxEvents(ctx, 100)
	if err != nil {
		log.Printf("Outbox failed to fetch events: %v", err)
		return
	}

	for _, event := range events {
		msg := kafka.Message{
			Key:   []byte(strconv.Itoa(int(event.AggregateID))),
			Value: event.Payload,
		}
		if err := r.writer.WriteMessages(ctx, msg); err != nil {
			log.Printf("Outbox failed to publish event %d: %v", event.ID, err)
		}
		if err := r.queries.MarkOutboxEventPublished(ctx, event.ID); err != nil {
			log.Printf("Outbox failed to mark event %d as published: %v", event.ID, err)
			return
		}
	}
}
