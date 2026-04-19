package main

import (
	"context"
	"log"
	"time"

	"github.com/1kyryll/go-grpc/internal/common/gen/kitchen"
	"github.com/1kyryll/go-grpc/internal/common/gen/orders"
	kitchenService "github.com/1kyryll/go-grpc/internal/services/kitchen/services"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func consumeOrderEvents(ctx context.Context, reader *kafka.Reader, svc *kitchenService.KitchenService, sse *SSEServer) {
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka consumer failed to fetch message: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var event orders.OrderEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Kafka consumer failed to unmarshal event: %v", err)
			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Kafka consumer failed to commit bad message: %v", err)
			}
			continue
		}

		log.Printf("NEW ORDER FROM Kafka: id=%d user=%d status=%s",
			event.OrderId, event.CustomerId, event.Status)

		// Create ticket idempotently via service
		_, err = svc.CreateTicket(ctx, &kitchen.Ticket{
			OrderId: event.OrderId,
			Status:  "OPEN",
		})
		if err != nil {
			log.Printf("Kafka consumer failed to create Ticket for Order %d: %v", event.OrderId, err)
			// Don't commit - will retry on next fetch
			continue
		}

		// Enrich order with item names via service
		items, err := svc.GetEnrichedOrderItems(ctx, event.OrderId)
		if err != nil {
			log.Printf("Kafka consumer failed to enrich Order %d: %v", event.OrderId, err)
			items = []string{}
		}

		sse.Broadcast(enrichedOrder{
			ID:     event.OrderId,
			UserID: event.CustomerId,
			Status: event.Status,
			Items:  items,
		})

		// Commit offset after successful processing
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Kafka consumer failed to commit offset: %v", err)
		}
	}
}
