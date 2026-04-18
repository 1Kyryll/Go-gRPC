package main

import (
	"context"
	"fmt"
	"log"

	"github.com/1kyryll/go-grpc/internal/common/gen/orders"
	"github.com/1kyryll/go-grpc/internal/common/sqlc"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func consumeOrderEvents(ctx context.Context, reader *kafka.Reader, queries *sqlc.Queries, sse *SSEServer) {
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka consumer failed to fetch message: %v", err)
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

		// Create ticket idempotently
		err = queries.CreateTicketsIdempotent(ctx, sqlc.CreateTicketsIdempotentParams{
			OrderID: event.OrderId,
			Status:  "OPEN",
		})
		if err != nil {
			log.Printf("Kafka consumer failed to create Ticket for Order %d: %v", event.OrderId, err)
			// Don't commit - will retry on next fetch
			continue
		}

		// Enrich the order with item names from the database
		enriched := enrichedOrder{
			ID:     event.OrderId,
			UserID: event.CustomerId,
			Status: event.Status,
			Items:  []string{},
		}

		orderItems, err := queries.GetOrderItemsByOrderIDs(ctx, []int32{event.OrderId})
		if err != nil {
			log.Printf("Kafka consumer failed to fetch items for Order %d: %v", event.OrderId, err)
		} else {
			menuIDs := make([]int32, len(orderItems))
			for i, oi := range orderItems {
				menuIDs[i] = oi.MenuItemID
			}

			menuItems, err := queries.GetMenuItemsByIDs(ctx, menuIDs)
			if err != nil {
				log.Printf("Kafka consumer failed to fetch Menu Items: %v", err)
			} else {
				nameMap := make(map[int32]string, len(menuItems))
				for _, mi := range menuItems {
					nameMap[mi.ID] = mi.Name
				}

				for _, oi := range orderItems {
					name := nameMap[oi.MenuItemID]
					if name == "" {
						name = fmt.Sprintf("Item #%d", oi.MenuItemID)
					}
					if oi.Quantity > 1 {
						enriched.Items = append(enriched.Items, fmt.Sprintf("%s x%d", name, oi.Quantity))
					} else {
						enriched.Items = append(enriched.Items, name)
					}
				}
			}
		}

		sse.Broadcast(enriched)

		// Commit offset after successful processing
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Kafka consumer failed to commit offset: %v", err)
		}
	}
}
