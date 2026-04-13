FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gateway ./cmd/gateway/

FROM alpine:latest

COPY --from=build /app/gateway /gateway

EXPOSE 8082

CMD ["/gateway"]
