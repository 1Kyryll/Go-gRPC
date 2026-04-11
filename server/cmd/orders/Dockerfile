FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o orders ./cmd/orders/

FROM alpine:latest

COPY --from=build /app/orders /orders

EXPOSE 8080
EXPOSE 9000

CMD ["/orders"]