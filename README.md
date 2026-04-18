# Restaurant System Microservices

A real-time restaurant order management system built with a microservices architecture. Orders flow from the customer-facing frontend through a GraphQL gateway, into gRPC-connected backend services, and out to a live kitchen dashboard — all orchestrated with Docker Compose.

## Architecture

```
┌──────────────┐     GraphQL/WS      ┌──────────────┐     gRPC stream      ┌──────────────┐
│              │  ─────────────────►  │              │  ─────────────────►  │              │
│    Client    │     :3000            │   Gateway    │     :8082            │    Orders    │
│  (Next.js)   │  ◄─────────────────  │  (GraphQL)   │  ◄─────────────────  │   Service    │
│              │                      │              │                      │  :8080/:9000 │
└──────────────┘                      └──────┬───────┘                      └──────┬───────┘
                                             │                                     │
                                        gRPC │                              outbox relay
                                             │                                     │
                                      ┌──────▼───────┐                      ┌──────▼───────┐
                                      │              │                      │              │
                                      │    User      │                      │    Kafka     │
                                      │   Service    │                      │  (KRaft)     │
                                      │ :8083/:9001  │                      │    :9092     │
                                      └──────────────┘                      └──────┬───────┘
                                                                                   │
                                                                            ┌──────▼───────┐
                                                                            │              │
                                                                     SSE    │   Kitchen    │
                                                                            │   Service    │
                                                                            │    :8081     │
                                                                            └──────┬───────┘
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
| **Gateway** | `:8082` | GraphQL API with queries, mutations, and WebSocket subscriptions. Aggregates data via dataloaders. Handles auth via JWT |
| **Orders** | `:8080` (HTTP) `:9000` (gRPC) | Core order service. Stores orders in PostgreSQL, writes events to a transactional outbox, and relays them to Kafka. Also broadcasts via gRPC streaming for Gateway subscriptions |
| **Kitchen** | `:8081` | Consumes order events from Kafka, creates kitchen tickets idempotently, and pushes enriched order data to the dashboard via SSE. Protected by JWT auth (kitchen staff only) |
| **User** | `:8083` (HTTP) `:9001` (gRPC) | Handles user registration, login, and JWT token generation with role-based claims |
| **Kafka** | `:9092` (internal) `:9094` (external) | Apache Kafka in KRaft mode. Topic `orders.events` carries order events keyed by `order_id` with Protobuf payloads |
| **Kafka UI** | `:8090` | Web dashboard for inspecting Kafka topics, messages, and consumer groups |
| **PostgreSQL** | `:5433` | Shared database for orders, users, menu items, tickets, and the outbox table |

### Communication Patterns

- **Client <-> Gateway**: GraphQL queries/mutations over HTTP, real-time status updates via WebSocket subscriptions (`graphql-transport-ws`)
- **Gateway <-> Orders**: gRPC for order creation; gRPC server streaming for `orderCreated` subscriptions
- **Gateway <-> User**: gRPC for registration, login, and user lookup
- **Orders -> Kafka**: Transactional outbox pattern — order events are written to an `outbox` table in the same DB transaction as the order, then a relay goroutine publishes them to Kafka topic `orders.events`
- **Kafka -> Kitchen**: Kitchen service consumes from `orders.events` (consumer group: `kitchen-service`), creates tickets idempotently (`ON CONFLICT DO NOTHING`), and enriches orders with menu item names
- **Kitchen -> Client**: Server-Sent Events (SSE) push enriched order data to the kitchen dashboard

## Authentication & Authorization

The system uses JWT-based authentication with role-based access control (RBAC).

### Roles

| Role | Capabilities |
|------|-------------|
| **CUSTOMER** | Browse menu, create orders, view/cancel own orders, track own order status |
| **KITCHEN_STAFF** | View all orders, update order statuses, complete orders, access kitchen dashboard and SSE stream |

### Auth Flow

1. User registers or logs in via GraphQL mutations (`register` / `login`)
2. The User service hashes the password (bcrypt), generates a JWT with claims (`sub`, `username`, `role`), and returns it
3. The frontend stores the token in a cookie and sends it as a `Bearer` token on GraphQL requests
4. The Gateway auth middleware extracts claims from the JWT and injects them into the request context
5. Each resolver checks the user's role before executing (server-side enforcement)
6. The Kitchen service validates JWTs independently — SSE uses `?access_token=` query param since `EventSource` cannot send headers

### Default Accounts

New registrations default to the `CUSTOMER` role. Kitchen staff accounts are seeded via migration or created manually in the database.

## Getting Started

### Prerequisites

Create a `.env` file in the project root:

```env
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=restaurant
JWT_SECRET=your_jwt_secret_here
```

### Docker Compose (recommended)

```bash
docker compose up -d
```

This starts all services. Visit:
- **Menu & ordering**: http://localhost:3000
- **Kitchen dashboard**: http://localhost:3000/kitchen
- **GraphQL playground**: http://localhost:8082
- **Kafka UI**: http://localhost:8090

### Local Development

#### Prerequisites

- Go 1.26+
- Node.js 20+
- PostgreSQL (or use Docker for just the DB: `docker compose up postgres -d`)
- Set environment variables: `GOOSE_DBSTRING`, `JWT_SECRET`

#### Server

```bash
cd server
go run ./cmd/orders    # gRPC :9000 + HTTP :8080
go run ./cmd/user      # gRPC :9001 + HTTP :8083
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
goose -dir internal/common/migrations postgres "$GOOSE_DBSTRING" up
```

#### Regenerate Protobuf / SQLC / GraphQL

```bash
cd server
make gen    # runs: make proto && make sqlc && make gqlgen
```

## API

### GraphQL (Gateway — `:8082`)

**Queries**
| Operation | Auth | Description |
|-----------|------|-------------|
| `order(id: ID!)` | Required (scoped) | Fetch a single order. Customers can only view their own |
| `orders(first: Int, status: OrderStatus)` | Required (scoped) | Paginated orders. Customers see only their own; kitchen staff sees all |
| `user(id: ID!)` | Required (scoped) | User details. Customers can only query themselves |
| `menuItems(first: Int, category: MenuCategory)` | None | Paginated menu with category filter |
| `search(query: String!)` | Optional (scoped) | Search across menu items, users, and orders. Results scoped by role |

**Mutations**
| Operation | Auth | Description |
|-----------|------|-------------|
| `register(input: RegisterInput!)` | None | Create a new customer account |
| `login(input: LoginInput!)` | None | Authenticate and receive JWT |
| `createOrder(input: CreateOrderInput!)` | Customer | Place a new order (must match authenticated user) |
| `updateOrderStatus(id: ID!, status: OrderStatus!)` | Kitchen staff | Change order status |
| `completeOrder(orderId: ID!)` | Kitchen staff | Mark order and ticket as completed |
| `cancelOrder(orderId: ID!)` | Customer | Cancel own order |

**Subscriptions**
| Operation | Auth | Description |
|-----------|------|-------------|
| `orderCreated` | Kitchen staff | Stream of newly created orders |
| `orderStatusChanged(orderId: ID)` | Required (scoped) | Real-time status updates. Customers must specify their own order ID |

### REST (Kitchen — `:8081`)

All endpoints require `KITCHEN_STAFF` role via JWT.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/stream?access_token=<jwt>` | SSE stream of new orders with item names |
| GET | `/ticket/{orderId}` | Get tickets for an order |
| POST | `/order/{orderId}/done` | Complete a ticket and its order |

### REST (Orders — `:8080`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/order/create` | Create a new order (also available via GraphQL) |

### REST (User — `:8083`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Login and receive JWT |

## Frontend Pages

| Route | Access | Description |
|-------|--------|-------------|
| `/` | Public | Menu with category filtering and cart controls |
| `/login` | Public | Login form |
| `/signup` | Public | Registration form |
| `/order/create` | Customer | Checkout with cart review and order submission |
| `/orders` | Customer | Own order list with status filtering |
| `/orders/[id]` | Customer | Order detail with live status via WebSocket |
| `/kitchen` | Kitchen staff | Real-time kitchen dashboard with SSE |
| `/ticket/[orderId]` | Kitchen staff | Ticket details for a specific order |

## Tech Stack

- **Go** — gRPC, protobuf, net/http, gqlgen (GraphQL), JWT (golang-jwt), bcrypt
- **PostgreSQL** — pgxpool, sqlc (type-safe queries), goose (migrations)
- **Next.js 16** — React 19, Tailwind CSS 4, Server Components, Server Actions
- **Messaging** — Apache Kafka (KRaft mode), transactional outbox pattern, `segmentio/kafka-go`
- **Real-time** — gRPC server streaming, GraphQL subscriptions (WebSocket), SSE
- **Auth** — JWT (HS256) with role-based access control, bcrypt password hashing
- **Infrastructure** — Docker multi-stage builds, Docker Compose, Kafka UI
