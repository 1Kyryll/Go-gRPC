# go-grpc

A real-time order management system built with Go, gRPC, and Next.js.

## Architecture

```
client/ → Next.js frontend (order forms, kitchen dashboard)
server/ → Go backend (orders service + kitchen service)
```

**Orders service** exposes both a gRPC server (`:9000`) and an HTTP server (`:8080`). Orders are stored in memory and broadcast to subscribers via gRPC server streaming.

**Kitchen service** connects to the orders gRPC stream and re-broadcasts new orders to browser clients over Server-Sent Events (`:8081`).

## Getting Started

### Server

```bash
cd server
go run ./services/orders    # starts gRPC :9000 + HTTP :8080
go run ./services/kitchen   # starts SSE :8081
```

### Client

```bash
cd client
npm install
npm run dev                 # starts Next.js on :3000
```

### Regenerate Protobuf

```bash
cd server
make gen
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/order/create` | Create a new order |
| GET | `/order/get` | List all orders |
| GET | `:8081/stream` | SSE stream of new orders |

## Tech Stack

- **Go** — gRPC, protobuf, net/http
- **Next.js** — React, Tailwind CSS, Server Actions
- **gRPC streaming** — real-time order broadcast
- **SSE** — push orders to the kitchen dashboard
