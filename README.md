# Restaurant System Microservices

A real-time restaurant order management system built with a microservices architecture. Orders flow from the customer-facing frontend through a GraphQL gateway, into gRPC-connected backend services, and out to a live kitchen dashboard — all orchestrated with Docker Compose.

## Architecture

```
┌──────────────┐     GraphQL/WS      ┌──────────────┐     gRPC stream      ┌──────────────┐
│              │  ─────────────────►  │              │  ─────────────────►  │              │
│    Client    │     :3000            │   Gateway    │     :8082            │    Orders    │
│  (Next.js)   │  ◄─────────────────  │  (GraphQL)   │  ◄─────────────────  │   Service    │
│              │                      │              │                      │  :8080/:9000 │
└──────────────┘                      └──────────────┘                      └──────┬───────┘
                                                                                   │
                                                                            gRPC stream
                                                                                   │
                                      ┌──────────────┐                      ┌──────▼───────┐
                                      │              │         SSE          │              │
                                      │   Kitchen    │  ◄─────────────────  │   Kitchen    │
                                      │  Dashboard   │     :8081            │   Service    │
                                      │  (via Client)│                      │              │
                                      └──────────────┘                      └──────────────┘
                                                                                   │
                                                                            ┌──────▼───────┐
                                                                            │  PostgreSQL  │
                                                                            │    :5433     │
                                                                            └──────────────┘
```

### Services

| Service | Port(s) | Description |
|---------|---------|-------------|
| **Client** | `:3000` | Next.js 16 frontend with menu browsing, cart, order placement, order tracking, and kitchen dashboard |
| **Gateway** | `:8082` | GraphQL API with queries, mutations, and WebSocket subscriptions. Aggregates data via dataloaders |
| **Orders** | `:8080` (HTTP) `:9000` (gRPC) | Core order service. Stores orders in PostgreSQL and broadcasts new orders via gRPC server streaming |
| **Kitchen** | `:8081` | Receives orders from the gRPC stream, creates kitchen tickets, and pushes enriched order data to the dashboard via SSE |
| **PostgreSQL** | `:5433` | Shared database for orders, customers, menu items, and tickets |

### Communication Patterns

- **Client ↔ Gateway**: GraphQL queries/mutations over HTTP, real-time status updates via WebSocket subscriptions (`graphql-transport-ws`)
- **Gateway ↔ Orders**: gRPC for order creation and streaming
- **Orders → Kitchen**: gRPC server streaming broadcasts new orders in real-time
- **Kitchen → Client**: Server-Sent Events (SSE) push enriched order data to the kitchen dashboard

## Getting Started

### Docker Compose (recommended)

```bash
docker compose up -d
```

This starts all five services. Visit:
- **Menu & ordering**: http://localhost:3000
- **Kitchen dashboard**: http://localhost:3000/kitchen
- **GraphQL playground**: http://localhost:8082

### Local Development

#### Prerequisites

- Go 1.26+
- Node.js 20+
- PostgreSQL (or use Docker for just the DB: `docker compose up postgres -d`)

#### Server

```bash
cd server
go run ./cmd/orders    # gRPC :9000 + HTTP :8080
go run ./cmd/kitchen   # HTTP + SSE :8081
go run ./cmd/gateway   # GraphQL :8082
```

#### Client

```bash
cd client
npm install
npm run dev            # Next.js on :3000
```

#### Database Migrations

```bash
cd server
goose -dir internal/services/common/migrations postgres "$GOOSE_DBSTRING" up
```

#### Regenerate Protobuf / SQLC / GraphQL

```bash
cd server
make gen       # protobuf
sqlc generate  # database queries
```

## API

### GraphQL (Gateway — `:8082`)

**Queries**
| Operation | Description |
|-----------|-------------|
| `order(id: ID!)` | Fetch a single order with items, customer, and ticket |
| `orders(first: Int, status: OrderStatus)` | Paginated order list with optional status filter |
| `customer(id: ID!)` | Customer details with their order history |
| `menuItems(first: Int, category: MenuCategory)` | Paginated menu with category filter |
| `search(query: String!)` | Full-text search across menu items, customers, and orders |

**Mutations**
| Operation | Description |
|-----------|-------------|
| `createOrder(input: CreateOrderInput!)` | Place a new order |
| `updateOrderStatus(id: ID!, status: OrderStatus!)` | Change order status |
| `completeOrder(orderId: ID!)` | Mark order and ticket as completed |
| `cancelOrder(orderId: ID!)` | Cancel an order |

**Subscriptions**
| Operation | Description |
|-----------|-------------|
| `orderCreated` | Stream of newly created orders |
| `orderStatusChanged(orderId: ID)` | Real-time status updates, optionally filtered by order |

### REST (Kitchen — `:8081`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/stream` | SSE stream of new orders with item names |
| GET | `/ticket/{orderId}` | Get tickets for an order |
| POST | `/order/{orderId}/done` | Complete a ticket and its order |

### REST (Orders — `:8080`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/order/create` | Create a new order (also available via GraphQL) |

## Frontend Pages

| Route | Description |
|-------|-------------|
| `/` | Menu with category filtering and cart controls |
| `/order/create` | Checkout with cart review and order submission |
| `/orders` | Order list with status filtering |
| `/orders/[id]` | Order detail with live status via WebSocket |
| `/kitchen` | Real-time kitchen dashboard with SSE |
| `/ticket/[orderId]` | Ticket details for a specific order |

## Tech Stack

- **Go** — gRPC, protobuf, net/http, gqlgen (GraphQL)
- **PostgreSQL** — pgxpool, sqlc (type-safe queries), goose (migrations)
- **Next.js 16** — React 19, Tailwind CSS 4, Server Components, Server Actions
- **Real-time** — gRPC server streaming, GraphQL subscriptions (WebSocket), SSE
- **Infrastructure** — Docker multi-stage builds, Docker Compose
