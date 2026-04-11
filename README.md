# Go-gRPC

A real-time order management system built with Go, gRPC, PostgreSQL, and Next.js.

## Architecture

```
client/  → Next.js frontend (order forms, ticket view, kitchen dashboard)
server/  → Go backend (orders service + kitchen service)
db/      → PostgreSQL (via Docker)
```

**Orders service** exposes a gRPC server (`:9000`) and an HTTP server (`:8080`). Orders are stored in PostgreSQL and broadcast to subscribers via gRPC server streaming.

**Kitchen service** connects to the orders gRPC stream, creates tickets in PostgreSQL, and re-broadcasts new orders to browser clients over Server-Sent Events (`:8081`). Handles ticket and order completion.

## Getting Started

### Docker Compose (recommended)

```bash
docker compose up -d
```

This starts all services:
- **Postgres** — `:5433`
- **Orders** — `:8080` (HTTP) + `:9000` (gRPC)
- **Kitchen** — `:8081` (HTTP + SSE)
- **Client** — `:3000` (Next.js)

### Local Development

#### Server

```bash
cd server
go run ./cmd/orders    # starts gRPC :9000 + HTTP :8080
go run ./cmd/kitchen   # starts SSE :8081
```

#### Client

```bash
cd client
npm install
npm run dev            # starts Next.js on :3000
```

#### Database Migrations

```bash
cd server
goose -dir internal/services/common/migrations postgres "$GOOSE_DBSTRING" up
```

### Regenerate Protobuf

```bash
cd server
make gen
```

## API

| Service | Method | Endpoint | Description |
|---------|--------|----------|-------------|
| Orders | POST | `/order/create` | Create a new order |
| Kitchen | GET | `/ticket/{orderId}` | Get tickets for an order |
| Kitchen | POST | `/order/{orderId}/done` | Complete a ticket and order |
| Kitchen | GET | `/stream` | SSE stream of new orders |

## Tech Stack

- **Go** — gRPC, protobuf, net/http
- **PostgreSQL** — pgxpool, sqlc, goose migrations
- **Next.js** — React, Tailwind CSS, Server Actions
- **gRPC streaming** — real-time order broadcast
- **SSE** — push orders to the kitchen dashboard
- **Docker** — multi-stage builds, Docker Compose
